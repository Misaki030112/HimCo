package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"web.misaki.world/FinalExam/aws"
	"web.misaki.world/FinalExam/xmly"
)

var (
	ConvertCount    int
	StoragePath     string
	TargetAlbumId   int
	OnlyConvert     bool
	DisableDownLoad bool
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	argsAnalyze()
	if TargetAlbumId == -1 {
		log.Fatalf("please specify the Album Id....")
	}
	pathDir := fmt.Sprintf("%s/%d", StoragePath, TargetAlbumId)
	audioDir := fmt.Sprintf("%s/%s", pathDir, "audio")
	convertPath := fmt.Sprintf("%s/%s", pathDir, "data")
	if OnlyConvert {
		for i := 1; i <= ConvertCount; i++ {
			s3Url, err := aws.UploadFile(fmt.Sprintf("%s/%d.mp4", audioDir, i))
			if err != nil {
				log.Printf("upload the audio file %s to aws error\n%v\n", fmt.Sprintf("%s/%d.mp4", audioDir, i), err)
				continue
			}
			aws.ConvertToText(s3Url, convertPath)
		}
		return
	}
	album := xmly.ObtainDetailForAlbumId(TargetAlbumId)
	if DisableDownLoad {
		album.WriteFile(pathDir)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		album.WriteFile(pathDir)
		wg.Done()
	}()
	go func() {
		album.DownLoadAudio(audioDir)
		log.Printf("the resource audio download complete")
		wg.Done()
	}()
	wg.Wait()
	// start convert to text
	for i := 1; i <= ConvertCount; i++ {
		s3Url, err := aws.UploadFile(fmt.Sprintf("%s/%d.mp4", audioDir, i))
		if err != nil {
			log.Printf("upload the audio file %s to aws error\n%v\n", fmt.Sprintf("%s/%d.mp4", audioDir, i), err)
			continue
		}
		aws.ConvertToText(s3Url, convertPath)
	}

}

func argsAnalyze() {
	flag.IntVar(&ConvertCount, "c", 0, "the audio to text convert count")
	flag.StringVar(&StoragePath, "s", "./", "the storage Parent Path")
	flag.IntVar(&TargetAlbumId, "t", -1, "the target Album Id")
	flag.BoolVar(&OnlyConvert, "n", false, "only convert to text")
	flag.BoolVar(&DisableDownLoad, "d", false, "disable download")
	flag.Parse()
}
