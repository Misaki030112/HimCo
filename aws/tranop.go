package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ConvertToText(mp4Url, savePath string) {
	jobName := mp4Url[strings.LastIndex(mp4Url, "/")+1:]
	_, err := TrClient.StartTranscriptionJob(context.TODO(), &transcribe.StartTranscriptionJobInput{
		Media: &types.Media{
			MediaFileUri: aws.String(mp4Url),
		},
		TranscriptionJobName:      aws.String(jobName),
		IdentifyMultipleLanguages: aws.Bool(true),
		LanguageOptions: []types.LanguageCode{
			types.LanguageCodeEnUs,
			types.LanguageCodeZhCn,
		},
		MediaFormat: types.MediaFormatMp4,
	})
	if err != nil {
		log.Panicf("for %s start the transcription job error,reason\n%v\n", mp4Url, err)
		return
	}

	for {
		jobOutPut, err := TrClient.GetTranscriptionJob(context.TODO(), &transcribe.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(jobName),
		})
		if err != nil {
			log.Panicf("try to get the TranscriptJob:[%s] fail,\n%v\n", jobName, err)
			return
		}
		if jobOutPut.TranscriptionJob.TranscriptionJobStatus == types.TranscriptionJobStatusFailed {
			log.Panicf("the TranscriptJob:[%s] execute fail,\nReason: %s", jobName, *jobOutPut.TranscriptionJob.FailureReason)
			return
		}
		if jobOutPut.TranscriptionJob.TranscriptionJobStatus == types.TranscriptionJobStatusCompleted {
			log.Printf("the TranscriptJob:[%s] execute success! completed Time:[%s]\n", jobName, jobOutPut.TranscriptionJob.CompletionTime.String())
			jobDetail := jobOutPut.TranscriptionJob
			if err := os.MkdirAll(savePath, 0750); err != nil {
				log.Panicf("can not create the file parent dir %s ....\n%v\n", savePath, err)
			}
			saveFilePath := filepath.Join(savePath, fmt.Sprintf("%s.json", jobName))
			f, err := os.Create(saveFilePath)
			if err != nil {
				log.Panicf("can not create the file %s .. \n%v\n", saveFilePath, err)
				return
			}
			res, err := http.Get(*jobDetail.Transcript.TranscriptFileUri)
			if err != nil {
				log.Panicf("request  url %s error \n%v\n", *jobDetail.Transcript.TranscriptFileUri, err)
				return
			}
			by, err := io.ReadAll(res.Body)
			if err != nil {
				log.Panicf("read the response data error\n%v\n", err)
				return
			}
			jsonObj := make(map[string]interface{})
			if err := json.Unmarshal(by, &jsonObj); err != nil {
				log.Panicf("the byte[] can not deserializes to jsonObj ...\n%v", err)
			}
			by, _ = json.MarshalIndent(jsonObj, " ", "  ")
			_, err = f.Write(by)
			if err != nil {
				log.Panicf("write the response data to the file %s error\n", saveFilePath)
			}
			err = f.Close()
			if err != nil {
				log.Printf("file %s close fail may case some content lost\n", saveFilePath)
				return
			}
			err = res.Body.Close()
			if err != nil {
				log.Printf("request %s response body close error \n", *jobDetail.Transcript.TranscriptFileUri)
				return
			}
			return
		}
		log.Printf("the TranscriptionJob:[%s] is executing......", jobName)
		time.Sleep(time.Duration(10) * time.Second)
	}

}
