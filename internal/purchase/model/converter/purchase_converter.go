package converter

import (
	"strconv"
	productDto "tutup-lapak/internal/product/dto"
	"tutup-lapak/internal/purchase/dto"
	"tutup-lapak/internal/purchase/model"
)

func ToPurchaseResponse(purchase model.Purchase, purchasedItems []productDto.ProductResponse, paymentDetails []dto.PaymentDetail) dto.PurchaseResponse {
	return dto.PurchaseResponse{
		PurchaseID:     strconv.Itoa(purchase.ID),
		TotalPrice:     purchase.TotalPrice,
		PurchasedItems: purchasedItems,
		PaymentDetails: paymentDetails,
	}
}
