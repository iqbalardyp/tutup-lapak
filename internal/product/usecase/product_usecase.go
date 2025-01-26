package usecase

import (
	"context"
	"tutup-lapak/internal/product/dto"
	"tutup-lapak/internal/product/repository"
)

type ProductUsecase struct {
	repo *repository.ProductRepo
}

func NewProductUsecase(repo *repository.ProductRepo) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
	}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, sellerID *int, payload *dto.ProductPayload) (*dto.ProductResponse, error) {
	product, err := u.repo.CreateProduct(ctx, sellerID, payload)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *ProductUsecase) GetProducts(ctx context.Context, payload *dto.ProductGetPayload) (*[]dto.ProductResponse, error) {
	products, err := u.repo.GetProducts(ctx, payload)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, ID, sellerID *int, payload *dto.ProductPayload) (*dto.ProductResponse, error) {
	product, err := u.repo.UpdateProduct(ctx, ID, sellerID, payload)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, ID, sellerID *int) error {
	err := u.repo.DeleteProduct(ctx, ID, sellerID)
	if err != nil {
		return err
	}
	return nil
}
