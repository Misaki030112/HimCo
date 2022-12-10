package xmly

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Album struct {
	Id              int      `json:"id"`
	Title           string   `json:"title,omitempty"`
	Mark            float32  `json:"mark,omitempty"`
	OriginPlayCount int64    `json:"-"`
	Subscribe       float32  `json:"subscribe,omitempty"`
	Labels          []string `json:"lable,omitempty"`
	Desc            string   `json:"desc,omitempty"`
	List            []Item   `json:"list,omitempty"`
}

type Item struct {
	Id              int     `json:"id"`
	Name            string  `json:"name,omitempty"`
	OriginPlayCount int64   `json:"-"`
	Subscribe       float32 `json:"subscribe,omitempty"`
	Date            string  `json:"date,omitempty"`
	HasAudio        bool    `json:"has-audio,omitempty"`
	AudioUrl        string  `json:"src,omitempty"`
}

func (album *Album) WriteFile(path string) {

	by, err := json.MarshalIndent(album, "", "  ")
	if err != nil {
		log.Panic("the album obj to json fail....", err)
	}
	if err := os.MkdirAll(path, 0750); err != nil {
		log.Panicf("can not create the file path: %s parent dir... : \n", path)
	}
	fileName := filepath.Join(path, "data.json")
	f, err := os.Create(fileName)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("the file %s close error,some content may lost", fileName)
		}
	}(f)
	if err != nil {
		log.Panicf("can not create the file , the file path: %s\n", path)
	}
	if _, err := f.WriteString(string(by)); err != nil {
		log.Panicf("the json string write file fail.. the json string : \n%s\n", string(by))
	}
	log.Printf("success write album to file: %s\n", fileName)
}

func (album *Album) DownLoadAudio(path string) {
	wg := sync.WaitGroup{}
	for l, r := 0, min(5, len(album.List)); r <= len(album.List); l, r = r, r+5 {
		wg.Add(1)
		l := l
		r := r
		go func() {
			album.batchDownload(path, l, min(r, len(album.List)))
			wg.Done()
		}()
	}
	wg.Wait()
}

// batchDownload batch MultiThread Download go execute func warn r is involved r<n
func (album *Album) batchDownload(path string, l int, r int) {
	for index, item := range album.List[l:r] {
		if !item.HasAudio {
			return
		}
		fileName := fmt.Sprintf("%d.mp4", l+index+1)
		fileName = filepath.Join(path, fileName)
		if err := os.MkdirAll(path, 0750); err != nil {
			log.Panicf("can not create the file parent dir %s ....\n%v\n", path, err)
		}
		f, err := os.Create(fileName)
		if err != nil {
			log.Printf("can not create the file %s .. \n%v\n", fileName, err)
			continue
		}
		res, err := http.Get(item.AudioUrl)
		if err != nil {
			log.Printf("request audio url %s error \n%v\n", item.AudioUrl, err)
			continue
		}
		_, err = io.Copy(f, res.Body)
		if err != nil {
			log.Printf("write the m4a data to the file %s error\n", fileName)
		}

		err = f.Close()
		if err != nil {
			log.Printf("file %s close fail may case some content lost\n", fileName)
		}
		err = res.Body.Close()
		if err != nil {
			log.Printf("request %s response body close error \n", item.AudioUrl)
		}
	}
}

func (album *Album) Top3Audio() []*Item {
	qItem := make(PriorityQueueItem, 0, 3)
	i := 0
	itemList := album.List
	for ; i < 2; i++ {
		qItem.Push(&itemList[i])
	}
	heap.Init(&qItem)
	for ; i < len(itemList); i++ {
		heap.Push(&qItem, &itemList[i])
		heap.Pop(&qItem)
	}
	res := make([]*Item, 3, 3)
	for i = 2; i >= 0; i-- {
		res[i] = heap.Pop(&qItem).(*Item)
	}
	return res
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

type PriorityQueueAlbum []*Album

func (pq PriorityQueueAlbum) Len() int {
	return len(pq)
}
func (pq PriorityQueueAlbum) Less(i, j int) bool {
	return pq[i].OriginPlayCount < pq[j].OriginPlayCount
}
func (pq PriorityQueueAlbum) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueueAlbum) Push(x any) {
	*pq = append(*pq, x.(*Album))
}
func (pq *PriorityQueueAlbum) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

type PriorityQueueItem []*Item

func (pq PriorityQueueItem) Len() int {
	return len(pq)
}

func (pq PriorityQueueItem) Less(i, j int) bool {
	return pq[i].OriginPlayCount < pq[j].OriginPlayCount
}

func (pq PriorityQueueItem) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueueItem) Push(x any) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueueItem) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
