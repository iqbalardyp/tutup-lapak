package handler

import (
	"mime/multipart"
	"net/http"
	file_usecase "tutup-lapak/internal/file/usecase"
	customErrors "tutup-lapak/pkg/custom-errors"
	"tutup-lapak/pkg/response"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FileHandler struct {
	Log         *logrus.Logger
	FileUsecase *file_usecase.FileUsecase
}

func NewFileHandler(Usecase *file_usecase.FileUsecase, logger *logrus.Logger) *FileHandler {
	return &FileHandler{
		Log:         logger,
		FileUsecase: Usecase,
	}
}

func (h *FileHandler) UploadFile(ctx echo.Context) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}
	defer file.Close()

	fileType, isValid := h.isValidFile(fileHeader, file)
	if !isValid {
		err = errors.Wrap(customErrors.ErrBadRequest, "file is invalid")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	fileResponse, err := h.FileUsecase.UploadFile(ctx.Request().Context(), file, *fileType)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &fileResponse)
}

func (h *FileHandler) isValidFile(fileHeader *multipart.FileHeader, file multipart.File) (*string, bool) {

	if fileHeader.Size > 100*1024 {
		return nil, false
	}

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return nil, false
	}
	// Reset the read pointer of the file
	if _, err := file.Seek(0, 0); err != nil {
		return nil, false
	}
	fileType := http.DetectContentType(buffer)

	switch fileType {
	case file_usecase.JPEG, file_usecase.JPG, file_usecase.PNG:
		return &fileType, true
	default:
		return nil, false
	}
}
