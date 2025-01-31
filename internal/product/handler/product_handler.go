package handler

import (
	"net/http"
	"strconv"
	"strings"
	"tutup-lapak/internal/product/dto"
	"tutup-lapak/internal/product/usecase"
	customErrors "tutup-lapak/pkg/custom-errors"
	"tutup-lapak/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type ProductHandler struct {
	usecase   *usecase.ProductUsecase
	validator *validator.Validate
}

var sortByCache = make(map[string]string)

const DEFAULT_LIMIT = 5

func NewProductHandler(usecase *usecase.ProductUsecase, validator *validator.Validate) *ProductHandler {
	return &ProductHandler{
		usecase:   usecase,
		validator: validator,
	}
}

func (h *ProductHandler) CreateProduct(ctx echo.Context) error {
	var payload dto.ProductPayload

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrBadRequest))
	}

	if err := h.validator.Struct(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	// TODO: Get sellerID from Auth
	sellerID := 1

	product, err := h.usecase.CreateProduct(ctx.Request().Context(), &sellerID, &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &product)
}

func (h *ProductHandler) GetProducts(ctx echo.Context) error {
	var payload dto.ProductGetPayload

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrBadRequest))
	}

	if err := h.validator.Struct(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if payload.Limit == 0 {
		payload.Limit = DEFAULT_LIMIT
	}

	if payload.SortBy != nil {
		sec, found := h.parseSortBy(payload.SortBy)
		if !found {
			return ctx.JSON(response.WriteErrorResponse(customErrors.ErrBadRequest))
		}

		payload.SortBy = sec
	}

	products, err := h.usecase.GetProducts(ctx.Request().Context(), &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if len(*products) == 0 {
		return ctx.JSON(http.StatusOK, make([]bool, 0))
	}

	return ctx.JSON(http.StatusOK, &products)
}

func (h *ProductHandler) UpdateProduct(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("productId"))
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrNotFound))
	}

	var payload dto.ProductPayload
	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrBadRequest))
	}

	if err := h.validator.Struct(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	// TODO: Get sellerID from Auth
	sellerID := 1

	product, err := h.usecase.UpdateProduct(ctx.Request().Context(), &id, &sellerID, &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &product)
}

func (h *ProductHandler) DeleteProduct(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("productId"))
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrNotFound))
	}

	// TODO: Get sellerID from Auth
	sellerID := 1

	err = h.usecase.DeleteProduct(ctx.Request().Context(), &id, &sellerID)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  "OK",
		Message: "Product is deleted",
	})
}

func (h *ProductHandler) parseSortBy(s *string) (*string, bool) {
	if s == nil {
		return nil, true
	}

	if !strings.Contains(*s, "sold-") {
		return s, true
	}

	if sec, found := sortByCache[*s]; found {
		return &sec, true
	}

	sec, found := strings.CutPrefix(*s, "sold-")
	if !found {
		return s, found
	}

	sortByCache[*s] = sec
	return &sec, true
}
