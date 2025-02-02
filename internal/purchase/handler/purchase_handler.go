package handler

import (
	"net/http"
	"strconv"

	"tutup-lapak/internal/purchase/dto"
	"tutup-lapak/internal/purchase/usecase"
	customErrors "tutup-lapak/pkg/custom-errors"
	"tutup-lapak/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type PurchaseHandler struct {
	UseCase  *usecase.PurchaseUseCase
	Validate *validator.Validate
}

func NewPurchaseHandler(useCase *usecase.PurchaseUseCase, validate *validator.Validate) *PurchaseHandler {
	return &PurchaseHandler{
		UseCase:  useCase,
		Validate: validate,
	}
}

func (h *PurchaseHandler) CreatePurchase(ctx echo.Context) error {
	var request = new(dto.PurchaseRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := h.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	purchase, err := h.UseCase.CreatePurchase(ctx.Request().Context(), request)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusCreated, purchase)
}

func (h *PurchaseHandler) CreatePayment(ctx echo.Context) error {
	purchaseIdStr := ctx.Param("purchaseId")
	purchaseId, err := strconv.Atoi(purchaseIdStr)
	if purchaseIdStr == "" || err != nil {
		err := errors.Wrap(customErrors.ErrNotFound, "purchase ID is required and must be a valid integer")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	var request = new(dto.PaymentRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := h.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	err = h.UseCase.CreatePayment(ctx.Request().Context(), purchaseId, request)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusCreated, response.BaseResponse{
		Status:  http.StatusText(http.StatusCreated),
		Message: "Successfully received payment",
	})
}
