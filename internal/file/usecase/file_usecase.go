package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"tutup-lapak/internal/file/dto"
	"tutup-lapak/internal/file/model/converter"
	"tutup-lapak/internal/file/repository"
	"tutup-lapak/pkg/dotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FileUsecase struct {
	S3Uploader *manager.Uploader
	Env        *dotenv.Env
	fileRepo   *repository.FileRepository
}

const (
	JPEG = "image/jpeg"
	JPG  = "image/jpg"
	PNG  = "image/png"
)

var (
	nameType = map[string]string{
		JPEG: ".jpeg",
		JPG:  ".jpg",
		PNG:  ".png",
	}
)

func NewFileUseCase(uploader *manager.Uploader, env *dotenv.Env, fileRepo *repository.FileRepository) *FileUsecase {
	return &FileUsecase{
		S3Uploader: uploader,
		Env:        env,
		fileRepo:   fileRepo,
	}
}

func (u *FileUsecase) UploadFile(ctx context.Context, file multipart.File, fileType string) (*dto.FileUploadResponse, error) {
	defer file.Close()

	filename := u.generateFilename(fileType)
	fileUri := u.generateFileUrl(filename)

	go func(uploader *manager.Uploader, file multipart.File, bucket, name string) {
		params := &s3.PutObjectInput{
			Bucket: aws.String(u.Env.AWS_S3_BUCKET_NAME),
			Key:    aws.String(filename),
			ACL:    types.ObjectCannedACLPublicRead,
			Body:   file,
		}
		_, err := uploader.Upload(context.Background(), params)
		if err != nil {
			fmt.Printf("failed to upload file: %v\n", err)
		}
	}(u.S3Uploader, file, u.Env.AWS_S3_BUCKET_NAME, filename)

	arg := repository.InsertFileParams{
		URI:          fileUri,
		ThumbnailURI: fileUri,
	}

	fileData, err := u.fileRepo.InsertFile(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upload file")
	}

	response := converter.ToFileResponse(fileData)
	return &response, nil
}

func (c *FileUsecase) generateFilename(fileType string) string {
	postfix := nameType[fileType]
	return uuid.New().String() + postfix
}

func (c *FileUsecase) generateFileUrl(filename string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		c.Env.AWS_S3_BUCKET_NAME,
		c.Env.AWS_S3_REGION,
		filename,
	)
}
