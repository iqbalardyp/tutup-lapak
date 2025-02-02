package dto

import (
	"tutup-lapak/internal/product/dto"
)

type ProductPurchaseRequest struct {
	ProductID string `json:"productId" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=1"`
}

type PurchaseRequest struct {
	PurchasedItems      []ProductPurchaseRequest `json:"purchasedItems" validate:"required,min=1,dive"`
	SenderName          string                   `json:"senderName" validate:"required,min=4,max=55"`
	SenderContactType   string                   `json:"senderContactType" validate:"required,oneof=email phone"`
	SenderContactDetail string                   `json:"senderContactDetail" validate:"required,contact_detail_validator"`
}

type PaymentRequest struct {
	FileIDs []string `json:"fileIds" validate:"required,min=1,dive"`
}

type PaymentDetail struct {
	SellerId          string `json:"sellerId"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
	TotalPrice        int    `json:"totalPrice"`
}

type PurchaseResponse struct {
	PurchaseID     string                `json:"purchaseId"`
	PurchasedItems []dto.ProductResponse `json:"purchasedItems"`
	TotalPrice     int                   `json:"totalPrice"`
	PaymentDetails []PaymentDetail       `json:"paymentDetails"`
}
