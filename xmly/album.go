package xmly

import (
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
	Title     string   `json:"title,omitempty"`
	Mark      float32  `json:"mark,omitempty"`
	Subscribe float32  `json:"subscribe,omitempty"`
	Labels    []string `json:"lable,omitempty"`
	Desc      string   `json:"desc,omitempty"`
	List      []Item   `json:"list,omitempty"`
}

type Item struct {
	Name      string  `json:"name,omitempty"`
	Subscribe float32 `json:"subscribe,omitempty"`
	Date      string  `json:"date,omitempty"`
	HasAudio  bool    `json:"-"`
	AudioUrl  string  `json:"-"`
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

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
