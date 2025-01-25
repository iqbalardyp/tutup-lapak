package config

import (
	"context"
	"log"
	"os"
	"tutup-lapak/pkg/dotenv"

	AWSConfig "github.com/aws/aws-sdk-go-v2/config"
	AWSCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader *manager.Uploader

var (
	AWS_S3_REGION      = os.Getenv("S3_REGION")
	AWS_S3_ID          = os.Getenv("S3_ID")
	AWS_S3_SECRET_KEY  = os.Getenv("S3_SECRET_KEY")
	AWS_S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME")
)

func NewS3Uploader(env *dotenv.Env) *manager.Uploader {
	config, err := AWSConfig.LoadDefaultConfig(
		context.TODO(),
		AWSConfig.WithRegion(env.AWS_S3_REGION),
		AWSConfig.WithCredentialsProvider(
			AWSCredentials.NewStaticCredentialsProvider(
				env.AWS_S3_ID,
				env.AWS_S3_SECRET_KEY,
				""),
		),
	)
	if err != nil {
		log.Fatal("unable connect to S3 Client", err.Error())
	}

	client := s3.NewFromConfig(config)
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // min size from aws
		u.Concurrency = 2            // vCPU max
		u.LeavePartsOnError = false
	})
	return uploader
}
