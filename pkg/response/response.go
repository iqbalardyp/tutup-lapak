package response

import (
	"net/http"
	customErrors "tutup-lapak/pkg/custom-errors"

	"github.com/pkg/errors"
)

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func WriteErrorResponse(err error) (int, BaseResponse) {
	cause := errors.Cause(err)
	msg := err.Error()

	switch cause {
	case customErrors.ErrNotFound:
		return http.StatusNotFound, BaseResponse{
			Status:  http.StatusText(http.StatusNotFound),
			Message: msg,
		}
	case customErrors.ErrConflict:
		return http.StatusConflict, BaseResponse{
			Status:  http.StatusText(http.StatusConflict),
			Message: msg,
		}
	case customErrors.ErrBadRequest:
		return http.StatusBadRequest, BaseResponse{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: msg,
		}
	case customErrors.ErrUnauthorized:
		return http.StatusUnauthorized, BaseResponse{
			Status:  http.StatusText(http.StatusUnauthorized),
			Message: msg,
		}
	default:
		return http.StatusInternalServerError, BaseResponse{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: msg,
		}
	}
}
