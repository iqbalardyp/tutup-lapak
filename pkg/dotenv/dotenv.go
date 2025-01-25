package dotenv

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	JWT_SECRET         string
	AWS_S3_REGION      string
	AWS_S3_ID          string
	AWS_S3_SECRET_KEY  string
	AWS_S3_BUCKET_NAME string
}

func LoadEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Env{
		JWT_SECRET:         os.Getenv("JWT_SECRET"),
		AWS_S3_REGION:      os.Getenv("S3_REGION"),
		AWS_S3_ID:          os.Getenv("S3_ID"),
		AWS_S3_SECRET_KEY:  os.Getenv("S3_SECRET_KEY"),
		AWS_S3_BUCKET_NAME: os.Getenv("S3_BUCKET_NAME"),
	}, nil
}
