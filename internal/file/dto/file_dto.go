package dto

type FileUploadResponse struct {
	FileID           string `json:"fileId"`
	FileURI          string `json:"fileUri"`
	FileThumbnailURI string `json:"fileThumbnailUri"`
}
