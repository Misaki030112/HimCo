package xmly

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	BasicInfoUrl    = "https://www.ximalaya.com/revision/album/v1/simple?albumId=%d"
	LabelsInfoUrl   = "https://www.ximalaya.com/revision/category/queryProductCategory?ptype=%d&id=%d&categoryId=%d"
	MarkInfoUrl     = "https://www.ximalaya.com/revision/comment/albumStatistics/%d"
	ItemTrackUrl    = "https://www.ximalaya.com/revision/album/v1/getTracksList?albumId=%d&pageNum=%d"
	ItemResourceUrl = "https://www.ximalaya.com/revision/play/v1/audio?id=%d&ptype=%d"
)

var client = &http.Client{}

// ObtainDetailForAlbumId  get the Album by Album-Id Example :
// https://www.ximalaya.com/album/2907021 the Id is 2907021
// warn: the func just return the album not ensure whether the album Items is ready,you can pass wg call wg.Wait()
func ObtainDetailForAlbumId(id int, wg *sync.WaitGroup) *Album {
	album := &Album{Id: id}
	album.Labels = make([]string, 0, 5)
	album.List = make([]*Item, 0, 300)
	//Get the Album ontology information
	res, err := EnhanceGet(fmt.Sprintf(BasicInfoUrl, id))
	if err != nil {
		wg.Done()
		log.Printf("%v", err)
		return nil
	}
	res = (res["data"].(map[string]interface{})["albumPageMainInfo"]).(map[string]interface{})
	album.Title = res["albumTitle"].(string)
	if res["shortIntro"] != nil {
		album.Desc = res["shortIntro"].(string)
	}
	album.Subscribe = float32(res["playCount"].(float64)) / 1e4
	album.OriginPlayCount = int64(res["playCount"].(float64))
	categoryId := int(res["categoryId"].(float64))

	res, err = EnhanceGet(fmt.Sprintf(LabelsInfoUrl, 1, id, categoryId))
	if err == nil {
		res = res["data"].(map[string]interface{})
		for _, m := range res["channels"].([]interface{}) {
			if m.(map[string]interface{})["channelName"] != nil {
				album.Labels = append(album.Labels, m.(map[string]interface{})["channelName"].(string))
			}
		}
	} else {
		log.Printf("%v", err)
	}

	res, err = EnhanceGet(fmt.Sprintf(MarkInfoUrl, id))
	if err == nil {
		res = res["data"].(map[string]interface{})
		if res["albumScore"] != nil {
			album.Mark = float32(res["albumScore"].(float64))
		}
	}
	pageNum := 1
	res, err = EnhanceGet(fmt.Sprintf(ItemTrackUrl, id, pageNum))
	if err != nil {
		log.Printf("can not get any Items of Album[id=%d] err:\n%v\n", id, err)
		if wg != nil {
			wg.Done()
		}
		return album
	}
	res = res["data"].(map[string]interface{})
	pageCount := (int(res["trackTotalCount"].(float64)) + int(res["pageSize"].(float64)) - 1) / int(res["pageSize"].(float64))
	if wg != nil {
		go func() {
			defer wg.Done()
			//Get the Album Items information
			getAlbumItems(pageNum, pageCount, res, err, id, album)

		}()
	} else {
		getAlbumItems(pageNum, pageCount, res, err, id, album)
	}
	return album
}

func getAlbumItems(pageNum int, pageCount int, res map[string]interface{}, err error, id int, album *Album) {
	for ; pageNum <= pageCount; pageNum++ {
		res, err = EnhanceGet(fmt.Sprintf(ItemTrackUrl, id, pageNum))
		if err != nil {
			log.Printf("can nor get this page's items pageNum=%d,pageCount=%d,Album[id=%d],err:\n%v\n", pageNum, pageCount, id)
			continue
		}
		res = res["data"].(map[string]interface{})
		for _, m := range res["tracks"].([]interface{}) {
			m := m.(map[string]interface{})
			trackId := int(m["trackId"].(float64))
			item := &Item{
				Subscribe:       float32(m["playCount"].(float64)) / 1e4,
				OriginPlayCount: int64(m["playCount"].(float64)),
				Id:              trackId,
			}
			if m["title"] != nil {
				item.Name = m["title"].(string)
			}
			if m["createDateFormat"] != nil {
				item.Date = m["createDateFormat"].(string)
			}
			resource, err := EnhanceGet(fmt.Sprintf(ItemResourceUrl, trackId, 1))
			if err != nil {
				log.Printf("can not get this Item[id=%d] information of Album[id=%d]..err:\n%v\n", trackId, id, err)
				continue
			}
			resource = (resource["data"]).(map[string]interface{})
			if resource["src"] != nil {
				item.HasAudio = true
				item.AudioUrl = resource["src"].(string)
			} else {
				item.HasAudio = false
				//log.Printf("the Album[%d] Item[%d]'s audio m4a file can not obtain,it has been encryption", id, trackId)
			}
			album.List = append(album.List, item)
		}
	}
}

// EnhanceGet Enhance the Get func and For anti-reptile mechanism
func EnhanceGet(url string) (map[string]interface{}, error) {
	res := Get(url)
	for i := 1; (res == nil || res["ret"] == nil || int(res["ret"].(float64)) != 200) && i < 50; i++ {
		//log.Printf("The Himalaya server rejected the request[%s] and will start reconnecting after 200ms, the %d retry", url, i)
		time.Sleep(200 * time.Microsecond)
		res = Get(url)
	}
	if res == nil || res["ret"] == nil || int(res["ret"].(float64)) != 200 {
		//log.Printf("request[%s] to Himalaya server error..", url)
		return nil, errors.New(fmt.Sprintf("can nor request to %s", url))
	}
	return res, nil
}

// Get send a Get request get json responseBody
func Get(url string) map[string]interface{} {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panicf("the request: %s is not correct ... \n%v", url, err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("the request: %s can not request ... \n%v", url, err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("the request: %s attem to close the response error, \n%v", url, err)
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("the request %s read response error \n%v", url, err)
		return nil
	}
	jsonObj := make(map[string]interface{})
	if err := json.Unmarshal(body, &jsonObj); err != nil {
		log.Printf("the request %s response body can not deserializes to jsonObj ...\n%v", url, err)
	}

	return jsonObj
}
