package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"web.misaki.world/HimCo/xmly"
)

// AnalyzeOutJson Analyze the Himalaya website data to get JSON data report
// Example /analyzeJson
func AnalyzeOutJson(w http.ResponseWriter, r *http.Request) {
	storageParentDir := r.Context().Value("StorageParentDir").(string)
	go func() {
		analyzeRes := xmly.AnalyzeXMLY()
		by, err := json.MarshalIndent(analyzeRes, "", "  ")
		if err != nil {
			log.Panicf("can not get The result after serialization analysis.....")
		}

		filePath := path.Join(storageParentDir, "analyze.json")
		if err := os.MkdirAll(storageParentDir, 0750); err != nil {
			log.Panicf("can not create the file path: %s parent dir... : \n", filePath)
		}
		f, err := os.Create(filePath)
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Printf("the file %s close error,some content may lost", filePath)
			}
		}(f)
		if err != nil {
			log.Panicf("can not create the file , the file path: %s\n", filePath)
		}
		if _, err := f.WriteString(string(by)); err != nil {
			log.Panicf("the json string write file fail.. the json string : \n%s\n", string(by))
		}
		log.Printf("Successfully analyzed data from the Himalayan website......\n")
	}()
	ReplyToText(w, "okk, The server starts to analyze the Himalayan website data .....")
}

func ReplyToJson(w http.ResponseWriter, obj any) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	jsonResp, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = io.WriteString(w, "sorry the server has been error....")
		if err != nil {
			log.Printf("response string [sorry the server has been error....] to client error\n%v\n", err)
			return
		}
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Printf("response json [%s] to client error\n%v\n", string(jsonResp), err)
	}
}
