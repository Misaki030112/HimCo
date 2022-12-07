package xmly

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
func ObtainDetailForAlbumId(id int) *Album {
	album := &Album{}
	album.Labels = make([]string, 0)
	album.List = make([]Item, 0)
	//Get the Album ontology information
	res := Get(fmt.Sprintf(BasicInfoUrl, id))
	res = (res["data"].(map[string]interface{})["albumPageMainInfo"]).(map[string]interface{})
	album.Title = res["albumTitle"].(string)
	album.Desc = res["shortIntro"].(string)
	album.Subscribe = float32(res["playCount"].(float64)) / 1e4

	categoryId := int(res["categoryId"].(float64))
	res = Get(fmt.Sprintf(LabelsInfoUrl, 1, id, categoryId))
	res = res["data"].(map[string]interface{})
	for _, m := range res["channels"].([]interface{}) {
		album.Labels = append(album.Labels, m.(map[string]interface{})["channelName"].(string))
	}

	res = Get(fmt.Sprintf(MarkInfoUrl, id))
	res = res["data"].(map[string]interface{})
	album.Mark = float32(res["albumScore"].(float64))

	//Get the Album Items information
	pageNum := 1
	res = Get(fmt.Sprintf(ItemTrackUrl, id, pageNum))
	res = res["data"].(map[string]interface{})
	pageCount := (int(res["trackTotalCount"].(float64)) + int(res["pageSize"].(float64)) - 1) / int(res["pageSize"].(float64))
	for ; pageNum <= pageCount; pageNum++ {
		res = (Get(fmt.Sprintf(ItemTrackUrl, id, pageNum))["data"]).(map[string]interface{})
		for _, m := range res["tracks"].([]interface{}) {
			m := m.(map[string]interface{})
			trackId := int(m["trackId"].(float64))
			item := Item{
				Name:      m["title"].(string),
				Subscribe: float32(m["playCount"].(float64)) / 1e4,
				Date:      m["createDateFormat"].(string),
			}
			resource := (Get(fmt.Sprintf(ItemResourceUrl, trackId, 1))["data"]).(map[string]interface{})
			if resource["src"] != nil {
				item.HasAudio = true
				item.AudioUrl = resource["src"].(string)
			} else {
				item.HasAudio = false
				log.Printf("the Album[%d] Item[%d]'s audio m4a file can not obtain,it has been encryption", id, trackId)
			}
			album.List = append(album.List, item)
		}
	}
	return album
}

// Get send a Get request get json responseBody
func Get(url string) map[string]interface{} {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panicf("the request: %s is not correct ... \n%v", url, err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Panicf("the request: %s can not request ... \n%v", url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("the request: %s attem to close the response error, \n%v", url, err)
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panicf("the request %s read response error \n%v", url, err)
	}
	jsonObj := make(map[string]interface{})
	if err := json.Unmarshal(body, &jsonObj); err != nil {
		log.Panicf("the request %s response body can not deserializes to jsonObj ...\n%v", url, err)
	}
	return jsonObj
}
