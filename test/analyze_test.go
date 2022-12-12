package test

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"testing"
	"web.misaki.world/FinalExam/xmly"
)

func TestAnalyzeXMLY(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzeRes := xmly.AnalyzeXMLY()
			by, err := json.MarshalIndent(analyzeRes, "", "  ")
			if err != nil {
				log.Panicf("can not get The result after serialization analysis.....")
			}

			filePath := path.Join("./", "analyze.json")
			if err := os.MkdirAll("./", 0750); err != nil {
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
		})
	}
}
