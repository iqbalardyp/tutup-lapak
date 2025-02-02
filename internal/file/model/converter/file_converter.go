package converter

import (
	"strconv"
	"tutup-lapak/internal/file/dto"
	"tutup-lapak/internal/file/model"
)

func ToFileResponse(file model.File) dto.FileUploadResponse {
	return dto.FileUploadResponse{
		FileID:           strconv.Itoa(file.ID),
		FileURI:          file.URI,
		FileThumbnailURI: file.ThumbnailURI,
	}
}
