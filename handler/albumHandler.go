package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"sync"

	"web.misaki.world/HimCo/aws"
	"web.misaki.world/HimCo/xmly"
)

// CrawlAlbum  deal with the request to CrawlAlbum
// Example  /album?id=284570&audioDownload=true
func CrawlAlbum(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		ReplyToText(w, "please input Query Param [id=?]....")
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 32)
	if err != nil {
		ReplyToText(w, "The input Query Param [id=?] is incorrect, it should be a 32-bit integer....")
	}
	audioDownload := false
	if r.URL.Query().Has("audioDownload") {
		audioDownload, _ = strconv.ParseBool(r.URL.Query().Get("audioDownload"))
	}
	storageParentDir := r.Context().Value("StorageParentDir").(string)

	go func() {
		log.Printf("start the fetching Task , fetch Album[id=%d]", id)
		pathDir := path.Join(storageParentDir, strconv.FormatInt(id, 10))
		wg := sync.WaitGroup{}
		wg.Add(1)
		album := xmly.ObtainDetailForAlbumId(int(id), &wg)
		wg.Wait()
		album.WriteFile(pathDir)
		if audioDownload {
			audioDir := path.Join(pathDir, "audio")
			album.DownLoadAudio(audioDir)
		}
		log.Printf("complete the fetching Task , fetch Album[id=%d]", id)
	}()
	ReplyToText(w, fmt.Sprintf("okk , The server starts fetching the Album[id=%d]", id))
}

// ConvertAudioToJson Convert audio to Json file and save
// Example /convert?id=284570&count=1
func ConvertAudioToJson(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		ReplyToText(w, "please input Query Param [id=?]....")
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 32)
	if err != nil {
		ReplyToText(w, "The input Query Param [id=?] is incorrect, it should be a 32-bit integer....")
	}

	count := int64(1)
	if r.URL.Query().Has("count") {
		count, _ = strconv.ParseInt(r.URL.Query().Get("count"), 10, 32)
	}
	storageParentDir := r.Context().Value("StorageParentDir").(string)

	go func() {
		log.Printf("start the audio convert Task , audio of Album[id=%d]", id)
		convertPath := path.Join(storageParentDir, strconv.FormatInt(id, 10), "data")
		audioDir := path.Join(storageParentDir, strconv.FormatInt(id, 10), "audio")
		for i := 1; i <= int(count); i++ {
			s3Url, err := aws.UploadFile(fmt.Sprintf("%s/%d.mp4", audioDir, i))
			if err != nil {
				log.Printf("upload the audio file %s to aws error\n%v\n", fmt.Sprintf("%s/%d.mp4", audioDir, i), err)
				continue
			}
			aws.ConvertToText(s3Url, convertPath)
		}
		log.Printf("complete the audio convert Task , audio of Album[id=%d]", id)
	}()
	ReplyToText(w, "okk , The server starts convert the audio of the Album[id=%d]")
}

func ReplyToText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := io.WriteString(w, s)
	if err != nil {
		log.Printf("response string [%s] to client error\n%v\n", s, err)
	}

}
