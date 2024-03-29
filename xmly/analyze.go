package xmly

import (
	"container/heap"
	"fmt"
	"log"
	"sync"
)

const (
	SubChannelsUrl = "https://www.ximalaya.com/revision/metadata/v2/group/channels?groupId=%d"
	AlbumsUrl      = "https://www.ximalaya.com/revision/metadata/v2/channel/albums?sort=%d&pageNum=%d&pageSize=%d&metadataValueId=%d"
)

// AnalyzeXMLY Analyze all Channels on the Himalayan website
func AnalyzeXMLY() []*Channel {
	channels := GetInitialChannels()
	analyzeRes := make([]*Channel, 0, 30)
	wg := sync.WaitGroup{}
	for id, channel := range channels {
		wg.Add(1)
		go analyzeChannel(id, channel, &wg)
	}
	wg.Wait()
	for _, channel := range channels {
		analyzeRes = append(analyzeRes, channel)
	}

	return analyzeRes
}

// analyzeChannel Analyze the parent Channel,And put task to analyze the child Channel
func analyzeChannel(id int, channel *Channel, wg *sync.WaitGroup) {
	defer wg.Done()
	res := Get(fmt.Sprintf(SubChannelsUrl, id))

	parentVipCount := int64(0)
	parentFinishCount := int64(0)
	allTotal := int64(0)

	channel.SubChannels = make([]*SubChannel, 0, 30)
	res = res["data"].(map[string]interface{})
	subChannelsJson := res["channels"].([]interface{})

	for _, subChannelJson := range subChannelsJson {
		wg.Add(1)

		m := subChannelJson.(map[string]interface{})
		metadataValueId := int64(m["relationMetadataValueId"].(float64))
		subChannel := &SubChannel{ChannelName: m["channelName"].(string)}
		res, err := EnhanceGet(fmt.Sprintf(AlbumsUrl, 3, 1, 50, metadataValueId))

		// Must ensure the first Page can get Albums
		// Enhanced processing for anti-crawling mechanism
		// Sometimes Himalaya will return you a status code of 200 but not a value for albums[].
		for err != nil || len((res["data"].(map[string]interface{}))["albums"].([]interface{})) == 0 {
			res, err = EnhanceGet(fmt.Sprintf(AlbumsUrl, 3, 1, 50, metadataValueId))
		}
		res = res["data"].(map[string]interface{})
		//only compute the first page audios of item It's enough
		albumsJson := res["albums"].([]interface{})
		albums := make([]*Album, 0, 50)
		computeWg := sync.WaitGroup{}

		for _, mJson := range albumsJson {
			m = mJson.(map[string]interface{})
			computeWg.Add(1)
			album := ObtainDetailForAlbumId(int(m["albumId"].(float64)), &computeWg)
			albums = append(albums, album)
		}

		total := int64(res["total"].(float64))
		subChannel.AlbumCount = total
		vipCount := int64(0)
		finishCount := int64(0)
		for i := 1; i <= (int(total)+99)/100; i++ {
			res, err = EnhanceGet(fmt.Sprintf(AlbumsUrl, 3, i, 100, metadataValueId))
			if res == nil {
				subChannel.AlbumCount -= 100
				continue
			}

			res = res["data"].(map[string]interface{})
			albumsJson = res["albums"].([]interface{})
			for _, mJson := range albumsJson {
				m = mJson.(map[string]interface{})
				if int64(m["vipType"].(float64)) != 0 {
					vipCount++
				}
				if int64(m["isFinished"].(float64)) == 1 {
					finishCount++
				}
			}
		}
		if total != 0 {
			subChannel.EndRate = int(finishCount * 100 / total)
			subChannel.VipRate = int(vipCount * 100 / total)
		}
		allTotal += total
		parentVipCount += vipCount
		parentFinishCount += finishCount

		go computeTop3(subChannel, albums, wg, &computeWg)
		albums = nil
		channel.SubChannels = append(channel.SubChannels, subChannel)
	}
	channel.SubChannelSize = len(channel.SubChannels)
	if allTotal != 0 {
		channel.VipRate = int(parentVipCount * 100 / allTotal)
		channel.EndRate = int(parentFinishCount * 100 / allTotal)
	}
}

// computeTop3 compute Node compute the top3 resources
func computeTop3(subChannel *SubChannel, albums []*Album, wg, computeWg *sync.WaitGroup) {
	defer wg.Done()
	computeWg.Wait()
	log.Printf("start compute Top3 of subChannel[ChannelName=%s]", subChannel.ChannelName)
	qAlbum := make(PriorityQueueAlbum, 0, 3)
	i := 0
	heap.Init(&qAlbum)
	for ; i < len(albums); i++ {
		heap.Push(&qAlbum, albums[i])
		if len(qAlbum) > 3 {
			heap.Pop(&qAlbum)
		}
	}

	subChannel.ShowTop3 = make([]*Album, 0, 3)
	for i = min(2, len(qAlbum)-1); i >= 0; i-- {
		subChannel.ShowTop3 = append(subChannel.ShowTop3, heap.Pop(&qAlbum).(*Album))
	}
	i = 0
	qItem := make(PriorityQueueItem, 0, 3)
	heap.Init(&qItem)
	for i = 0; i < len(albums); i++ {
		for _, v := range (albums[i]).Top3Audio() {
			heap.Push(&qItem, v)
			if len(qItem) > 3 {
				heap.Pop(&qItem)
			}
		}
	}
	subChannel.AudioTop3 = make([]*Item, 0, 3)
	for i = min(2, len(qItem)-1); i >= 0; i-- {
		subChannel.AudioTop3 = append(subChannel.AudioTop3, heap.Pop(&qItem).(*Item))
	}
	for _, v := range subChannel.ShowTop3 {
		v.List = nil
	}
	log.Printf("end compute Top3 of subChannel[ChannelName=%s]", subChannel.ChannelName)
}
