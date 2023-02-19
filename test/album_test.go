package test

import (
	"fmt"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"web.misaki.world/HimCo/xmly"
)

func TestAlbum_DownLoadAudio(t *testing.T) {
	type fields struct {
		Title     string
		Mark      float32
		Subscribe float32
		Labels    []string
		Desc      string
		List      []*xmly.Item
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			fields: fields{
				List: []*xmly.Item{
					{
						HasAudio: true,
						AudioUrl: "https://aod.cos.tx.xmcdn.com/group54/M07/57/F0/wKgLclwxLiShNNynAA_4a_o7K-k365.m4a",
					},
				},
			},
			args: args{
				path: "20403413/audio",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			album := &xmly.Album{
				Title:     tt.fields.Title,
				Mark:      tt.fields.Mark,
				Subscribe: tt.fields.Subscribe,
				Labels:    tt.fields.Labels,
				Desc:      tt.fields.Desc,
				List:      tt.fields.List,
			}
			album.DownLoadAudio(tt.args.path)
		})
	}
}

func TestAlbum_WriteFile(t *testing.T) {
	type fields struct {
		Title     string
		Mark      float32
		Subscribe float32
		Labels    []string
		Desc      string
		List      []*xmly.Item
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			fields: fields{
				Title:     "跟Lily说英语去旅行",
				Mark:      9.3,
				Subscribe: 207.9,
				Labels:    []string{"口语", "英语", "上班族", "听力"},
				Desc:      "跟Lily一起说英语去旅行的训练营即将开营啰！ 有144节线上课程，针对24个不同的旅游场景循环加深强度，课后你还可以缴交自己的录音还有老师亲自帮助你纠正不好的发音，让你立即开口说英语，在",
				List: []*xmly.Item{
					{
						Name:      "Lesson 1：在机场 At the Airport",
						Subscribe: 16.9,
						Date:      "2019-01",
					},
					{
						Name:      "Lesson 2：在飞机上 On the Plane",
						Subscribe: 11.6,
						Date:      "2019-01",
					},
				},
			},
			args: args{
				path: "20403413",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			album := &xmly.Album{
				Title:     tt.fields.Title,
				Mark:      tt.fields.Mark,
				Subscribe: tt.fields.Subscribe,
				Labels:    tt.fields.Labels,
				Desc:      tt.fields.Desc,
				List:      tt.fields.List,
			}
			album.WriteFile(tt.args.path)
		})
	}
}

func TestStatic_Album100W(t *testing.T) {
	var sum int32 = 0
	channels := xmly.GetInitialChannels()
	allWg := sync.WaitGroup{}
	for id, _ := range channels {
		res := xmly.Get(fmt.Sprintf(xmly.SubChannelsUrl, id))
		res = res["data"].(map[string]interface{})
		subChannelsJson := res["channels"].([]interface{})
		for _, subChannelJson := range subChannelsJson {
			allWg.Add(1)
			go func() {
				wg := sync.WaitGroup{}
				albums := make([]*xmly.Album, 0, 3005)
				m := subChannelJson.(map[string]interface{})
				metadataValueId := int64(m["relationMetadataValueId"].(float64))
				res1, err := xmly.EnhanceGet(fmt.Sprintf(xmly.AlbumsUrl, 2, 1, 1, metadataValueId))
				for err != nil || len((res1["data"].(map[string]interface{}))["albums"].([]interface{})) == 0 {
					res1, err = xmly.EnhanceGet(fmt.Sprintf(xmly.AlbumsUrl, 2, 1, 1, metadataValueId))
				}
				res1 = res1["data"].(map[string]interface{})
				total := int64(res1["total"].(float64))
				page := int(total) / 100

				for i := 1; i <= page; i++ {
					res1, err = xmly.EnhanceGet(fmt.Sprintf(xmly.AlbumsUrl, 2, i, 100, metadataValueId))
					for err != nil || len((res1["data"].(map[string]interface{}))["albums"].([]interface{})) == 0 {
						res1, err = xmly.EnhanceGet(fmt.Sprintf(xmly.AlbumsUrl, 2, i, 100, metadataValueId))
					}
					res1 = res1["data"].(map[string]interface{})
					albumsJson := res1["albums"].([]interface{})
					for _, mJson := range albumsJson {
						m1 := mJson.(map[string]interface{})
						AlbumId := int(m1["albumId"].(float64))
						if atomic.AddInt32(&sum, 1) > 1e6 {
							wg.Wait()
							for _, albumFile := range albums {
								albumFile.WriteFile(path.Join("../books", strconv.FormatInt(int64(albumFile.Id), 10)))
							}
							allWg.Done()
							return
						}
						t.Logf("currently fetched  %d book", sum)
						wg.Add(1)
						album := xmly.ObtainDetailForAlbumId(AlbumId, &wg)
						albums = append(albums, album)
					}
				}
				wg.Wait()
				for _, albumFile := range albums {
					albumFile.WriteFile(path.Join("../books", strconv.FormatInt(int64(albumFile.Id), 10)))
				}
				allWg.Done()
			}()
		}
		allWg.Wait()
	}

}
