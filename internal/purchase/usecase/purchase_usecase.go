package usecase

import (
	"context"
	"strconv"

	productDto "tutup-lapak/internal/product/dto"
	productRepository "tutup-lapak/internal/product/repository"
	"tutup-lapak/internal/purchase/dto"
	"tutup-lapak/internal/purchase/model/converter"
	"tutup-lapak/internal/purchase/repository"
	customErrors "tutup-lapak/pkg/custom-errors"
	"tutup-lapak/pkg/helper"

	"github.com/pkg/errors"
)

type PurchaseUseCase struct {
	purchaseRepo *repository.PurchaseRepository
	productRepo  *productRepository.ProductRepo
}

func NewPurchaseUseCase(purchaseRepo *repository.PurchaseRepository, productRepo *productRepository.ProductRepo) *PurchaseUseCase {
	return &PurchaseUseCase{
		purchaseRepo,
		productRepo,
	}
}

func (u *PurchaseUseCase) CreatePurchase(ctx context.Context, request *dto.PurchaseRequest) (*dto.PurchaseResponse, error) {
	var productIDs []int
	var purchasedQuantityMap = make(map[string]int)
	var paymentDetailsMap = make(map[string]dto.PaymentDetail)
	var purchasedItems []productDto.ProductResponse
	totalPrice := 0

	for _, item := range request.PurchasedItems {
		productID, err := strconv.Atoi(item.ProductID)
		if err != nil {
			return nil, errors.Wrap(customErrors.ErrBadRequest, "invalid product ID")
		}
		productIDs = append(productIDs, productID)
		purchasedQuantityMap[item.ProductID] = item.Qty
	}

	productData, err := u.productRepo.GetProductsByIDs(ctx, productIDs)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get products")
	}

	if len(productData) == 0 {
		return nil, errors.Wrap(customErrors.ErrBadRequest, "failed to get products")
	}

	for _, item := range productData {
		purchasedQuantity := purchasedQuantityMap[item.ProductID]

		if purchasedQuantity > item.Qty {
			return nil, errors.Wrap(customErrors.ErrBadRequest, "The requested quantity is unavailable or exceeds the available stock.")
		}

		purchasedItem := item.ProductResponse
		purchasedItems = append(purchasedItems, purchasedItem)

		paymentDetail := dto.PaymentDetail{
			SellerId:          item.SellerId,
			BankAccountName:   item.BankAccountName,
			BankAccountHolder: item.BankAccountHolder,
			BankAccountNumber: item.BankAccountNumber,
			TotalPrice:        item.Price * purchasedQuantity,
		}
		if detail, exists := paymentDetailsMap[item.SellerId]; exists {
			detail.TotalPrice += paymentDetail.TotalPrice
			paymentDetailsMap[item.SellerId] = detail // Update map
		} else {
			paymentDetailsMap[item.SellerId] = paymentDetail
		}

		totalPrice += paymentDetail.TotalPrice
	}
	paymentDetails := helper.MapToSlice(paymentDetailsMap)

	arg := repository.CreatePurchaseParams{
		TotalPrice:          totalPrice,
		TotalTransfer:       len(paymentDetails),
		SenderName:          request.SenderName,
		SenderContactType:   request.SenderContactType,
		SenderContactDetail: request.SenderContactDetail,
		PurchasedItems:      request.PurchasedItems,
	}

	purchase, err := u.purchaseRepo.CreatePurchase(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create purchase")
	}

	response := converter.ToPurchaseResponse(purchase, purchasedItems, paymentDetails)
	return &response, nil
}

func (u *PurchaseUseCase) CreatePayment(ctx context.Context, purchaseId int, request *dto.PaymentRequest) error {
	var fileIDs []int
	for _, item := range request.FileIDs {
		fileID, err := strconv.Atoi(item)
		if err != nil {
			return errors.Wrap(customErrors.ErrBadRequest, "invalid file ID")
		}
		fileIDs = append(fileIDs, fileID)
	}

	purchase, err := u.purchaseRepo.GetPurchase(ctx, purchaseId)
	if err != nil {
		return errors.Wrap(err, "failed to get purchase")
	}

	purchaseProducts, err := u.purchaseRepo.GetPurchaseProductsById(ctx, purchaseId)
	if err != nil {
		return errors.Wrap(err, "failed to get products")
	}

	if len(fileIDs) != purchase.TotalTransfer {
		return errors.Wrap(customErrors.ErrBadRequest, "missing payment")
	}

	arg := repository.UpdatePurchaseParams{
		PurchaseID:       purchaseId,
		PurchaseProducts: purchaseProducts,
	}

	err = u.purchaseRepo.UpdatePurchase(ctx, arg)
	if err != nil {
		return errors.Wrap(err, "failed to receive payment")
	}

	return nil
}
