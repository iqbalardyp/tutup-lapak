package repository

import (
	"context"
	"fmt"
	"tutup-lapak/internal/product/dto"
	customErrors "tutup-lapak/pkg/custom-errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

const (
	queryCreateProduct = `
	WITH product as (
		INSERT INTO products (seller_id, name, category, qty, price, sku, file_id)
		VALUES (@sellerID, @name, @category, @qty, @price, @sku, @fileID)
		RETURNING id::TEXT id, name, category, qty, price, sku, created_at, updated_at, file_id
	)
	SELECT
		p.id,
		p.name,
		p.category,
		p.qty,
		p.price,
		p.sku,
		p.created_at,
		p.updated_at,
		f.id::TEXT file_id,
		f.uri file_uri,
		f.thumbnail_uri file_thumbnail_uri
	FROM product p
	JOIN files f ON f.id = p.file_id;`
	queryUpdateProduct = `
	WITH product as (
		UPDATE products 
		SET 
			name = @name,
			category = @category,
			qty = @qty,
			price = @price,
			sku = @sku,
			file_id = @fileID::BIGINT
		WHERE
			id = @ID::BIGINT AND seller_id = @sellerID
		RETURNING id::TEXT id, name, category, qty, price, sku, updated_at, created_at, file_id
	)
	SELECT
		p.id,
		p.name,
		p.category,
		p.qty,
		p.price,
		p.sku,
		p.created_at,
		p.updated_at,
		f.id::TEXT file_id,
		f.uri file_uri,
		f.thumbnail_uri file_thumbnail_uri
	FROM product p
	JOIN files f ON f.id = p.file_id;`
	queryDeleteProductFromPivot    = "DELETE FROM pivot_purchase_products WHERE product_id = @ID;"
	queryDeleteProductFromProducts = "DELETE FROM products WHERE id = @ID AND seller_id = @sellerID;"
	queryGetProducts               = `
	SELECT
		p.id::TEXT id,
		p.name,
		p.category,
		p.qty,
		p.price,
		p.sku,
		p.updated_at,
		p.created_at,
		f.id::TEXT file_id,
		f.uri file_uri,
		f.thumbnail_uri file_thumbnail_uri
	FROM products p
	JOIN files f ON f.id = p.file_id
	WHERE
		(@productID::BIGINT IS NULL OR p.id = @productID::BIGINT)
		AND (@sku::TEXT IS NULL OR p.sku = @sku::TEXT)
		AND (@category::enum_product_categories IS NULL OR p.category = @category::enum_product_categories)
		AND (COALESCE(@sortBy::TEXT, '') !~ '^[0-9]+$'
			OR (COALESCE(@sortBy::TEXT, '') ~ '^[0-9]+$' AND p.id IN (
				SELECT DISTINCT ON (ppp.product_id) ppp.product_id
				FROM pivot_purchase_products ppp
				WHERE ppp.created_at >= NOW() - (@sortBy::TEXT || ' seconds')::INTERVAL
				ORDER BY ppp.product_id, ppp.created_at DESC
				LIMIT @limit
				OFFSET @offset
				))
			)
	ORDER BY
		CASE 
			WHEN @sortBy::TEXT = 'newest' THEN GREATEST(p.created_at, p.updated_at)
		END DESC,
		CASE 
			WHEN @sortBy::TEXT = 'cheapest' THEN p.price 
		END ASC
	LIMIT @limit
	OFFSET @offset;`
)

func (r *ProductRepo) CreateProduct(ctx context.Context, sellerID *int, payload *dto.ProductPayload) (*dto.ProductResponse, error) {
	var product dto.ProductResponse
	args := pgx.NamedArgs{
		"sellerID": &sellerID,
		"name":     &payload.Name,
		"category": &payload.Category,
		"qty":      &payload.Qty,
		"price":    &payload.Price,
		"sku":      &payload.Sku,
		"fileID":   &payload.FileID,
	}

	err := r.db.QueryRow(ctx, queryCreateProduct, args).Scan(
		&product.ProductID,
		&product.Name,
		&product.Category,
		&product.Qty,
		&product.Price,
		&product.Sku,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.FileID,
		&product.FileURI,
		&product.FileThumbnailURI,
	)
	if err != nil {
		return nil, customErrors.HandlePgError(err, "failed create product")
	}

	return &product, nil
}

func (r *ProductRepo) GetProducts(ctx context.Context, payload *dto.ProductGetPayload) (*[]dto.ProductResponse, error) {
	args := pgx.NamedArgs{
		"limit":     &payload.Limit,
		"offset":    &payload.Offset,
		"productID": &payload.ProductID,
		"category":  &payload.Category,
		"sku":       &payload.Sku,
		"sortBy":    &payload.SortBy,
	}

	rows, err := r.db.Query(ctx, queryGetProducts, args)
	if err != nil {
		return nil, customErrors.HandlePgError(err, "failed get product")
	}
	defer rows.Close()

	var products []dto.ProductResponse
	for rows.Next() {
		var product dto.ProductResponse
		if err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Category,
			&product.Qty,
			&product.Price,
			&product.Sku,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.FileID,
			&product.FileURI,
			&product.FileThumbnailURI,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return &products, nil
}

func (r *ProductRepo) UpdateProduct(ctx context.Context, ID, sellerID *int, payload *dto.ProductPayload) (*dto.ProductResponse, error) {
	var product dto.ProductResponse
	args := pgx.NamedArgs{
		"ID":       &ID,
		"sellerID": &sellerID,
		"name":     &payload.Name,
		"category": &payload.Category,
		"qty":      &payload.Qty,
		"price":    &payload.Price,
		"sku":      &payload.Sku,
		"fileID":   &payload.FileID,
	}

	err := r.db.QueryRow(ctx, queryUpdateProduct, args).Scan(
		&product.ProductID,
		&product.Name,
		&product.Category,
		&product.Qty,
		&product.Price,
		&product.Sku,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.FileID,
		&product.FileURI,
		&product.FileThumbnailURI,
	)
	if err != nil {
		return nil, customErrors.HandlePgError(err, "failed update product")
	}

	return &product, nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, ID, sellerID *int) error {
	deleteFromPivotArgs := pgx.NamedArgs{"ID": &ID}
	deleteFromProductsArgs := pgx.NamedArgs{"ID": &ID, "sellerID": &sellerID}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return customErrors.HandlePgError(err, "could not begin transaction")
	}
	batch := &pgx.Batch{}

	batch.Queue(queryDeleteProductFromPivot, deleteFromPivotArgs)
	batch.Queue(queryDeleteProductFromProducts, deleteFromProductsArgs)

	batchResults := tx.SendBatch(ctx, batch)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Delete product from pivot
	_, err = batchResults.Exec()
	if err != nil {
		return customErrors.HandlePgError(err, "could not execute batch query")
	}

	// Delete product from products
	result, err := batchResults.Exec()
	if err != nil {
		return customErrors.HandlePgError(err, "could not execute batch query")
	}
	if result.RowsAffected() != 1 {
		err = customErrors.ErrNotFound
		return customErrors.HandlePgError(err, "product not found")
	}

	if err := batchResults.Close(); err != nil {
		return fmt.Errorf("could not close batch result: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}
	return nil
}
