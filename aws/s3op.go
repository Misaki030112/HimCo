package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

const BucketName = "audio-2022-12-08"

// UploadFile upload the file to AWS S3 bucket return the resource url
func UploadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Panicf("could't open the file %s to upload, Here is why:\n %v\n", filePath, err)
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("close the file error ,may cause the file lost some content...\n%v\n", err)
		}
	}(file)
	_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filePath),
		Body:   file,
	})
	if err != nil {
		log.Panicf("Couldn't upload file %s to %s:%s. Here's why:\n %v\n",
			filePath, BucketName, filePath, err)
		return "", err
	}
	return fmt.Sprintf("s3://%s/%s", BucketName, filePath), nil
}
