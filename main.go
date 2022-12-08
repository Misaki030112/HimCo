package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"web.misaki.world/FinalExam/aws"
	"web.misaki.world/FinalExam/handler"
	"web.misaki.world/FinalExam/xmly"
)

var (
	ConvertCount          int
	StoragePath           string
	TargetAlbumId         int
	OnlyConvert           bool
	DisableDownLoad       bool
	TargetAlbumIdFilePath string
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	mux := http.NewServeMux()
	mux.HandleFunc("/album", handler.CrawlAlbum)

	err := http.ListenAndServe(":80", mux)
	log.Panicf("can not start server,Here is the reason:\n%v\n", err)
}

func main1() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	argsAnalyze()
	if TargetAlbumId == -1 && TargetAlbumIdFilePath == "" {
		log.Fatalf("please specify the Album Id.... or specify the Album Id File")
	}
	if TargetAlbumIdFilePath != "" {
		f, err := os.OpenFile(TargetAlbumIdFilePath, os.O_RDONLY, 0222)
		if err != nil {
			log.Fatalf("open the file %s ,error:\n %v\n", TargetAlbumIdFilePath, err)
		}
		fileScanner := bufio.NewScanner(f)
		for fileScanner.Scan() {
			id, err := strconv.ParseInt(fileScanner.Text(), 10, 32)
			if err != nil {
				log.Fatalf("Content encountered in file that cannot be converted to a number，error:\n%v\n", err)
			}
			doHandler(int(id))
		}
		return
	}
	doHandler(TargetAlbumId)

}

func doHandler(id int) {
	pathDir := fmt.Sprintf("%s/%d", StoragePath, id)
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
	album := xmly.ObtainDetailForAlbumId(id)
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
	flag.StringVar(&TargetAlbumIdFilePath, "f", "", "the target Album Id File")
	flag.Parse()
}
