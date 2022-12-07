package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"log"
)

var (
	S3Client *s3.Client
	TrClient *transcribe.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("error,the aws config not correct:\n%v\n", err)
	}
	S3Client = s3.NewFromConfig(cfg)
	TrClient = transcribe.NewFromConfig(cfg)
}
