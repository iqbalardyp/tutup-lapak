package dto

import "time"

type ProductPayload struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	Category string `json:"category" validate:"required,oneof=Food Beverage Clothes Furniture Tools"`
	Qty      int    `json:"qty" validate:"required,number,min=1"`
	Price    int    `json:"price" validate:"required,number,min=100"`
	Sku      string `json:"sku" validate:"required,min=1,max=32"`
	FileID   string `json:"fileId" validate:"required,number"`
}

type ProductGetPayload struct {
	Limit     int     `query:"limit" validate:"omitempty,number,min=0"`
	Offset    int     `query:"offset" validate:"omitempty,number,min=0"`
	ProductID *string `query:"productId" validate:"omitempty,number,min=1"`
	Sku       *string `query:"sku" validate:"omitempty,min=1"`
	Category  *string `query:"category" validate:"omitempty,oneof=Food Beverage Clothes Furniture Tools"`
	SortBy    *string `query:"sortBy" validate:"omitempty,sort_by"`
}

type ProductResponse struct {
	ProductID        string    `json:"productId"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Qty              int       `json:"qty"`
	Price            int       `json:"price"`
	Sku              string    `json:"sku"`
	FileID           string    `json:"fileId"`
	FileURI          string    `json:"fileUri"`
	FileThumbnailURI string    `json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type ProductWithSeller struct {
	ProductResponse
	SellerId          string
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
}
