package repository

import (
	"context"
	"time"

	"tutup-lapak/internal/purchase/dto"
	"tutup-lapak/internal/purchase/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepository struct {
	pool *pgxpool.Pool
}

func NewPurchaseRepository(pool *pgxpool.Pool) *PurchaseRepository {
	return &PurchaseRepository{pool: pool}
}

const createPurchaseQuery = `-- name: CreatePurchase :one
INSERT INTO purchases (
  total_price, total_transfer, sender_name, sender_contact_type, sender_contact_detail, paid_at
) VALUES (
  $1, $2, $3, $4, $5, NULL
) RETURNING id, total_price, total_transfer, sender_name, sender_contact_type, sender_contact_detail, paid_at
`

const insertPurchaseProductsQuery = `-- name: InsertPurchaseProducts :exec
INSERT INTO pivot_purchase_products (purchase_id, product_id, qty) VALUES ($1, $2, $3)
`

type CreatePurchaseParams struct {
	TotalPrice          int
	TotalTransfer       int
	SenderName          string
	SenderContactType   string
	SenderContactDetail string
	PurchasedItems      []dto.ProductPurchaseRequest
}

func (r *PurchaseRepository) CreatePurchase(ctx context.Context, arg CreatePurchaseParams) (model.Purchase, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.Purchase{}, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, createPurchaseQuery,
		arg.TotalPrice,
		arg.TotalTransfer,
		arg.SenderName,
		arg.SenderContactType,
		arg.SenderContactDetail,
	)

	var purchase model.Purchase
	err = row.Scan(
		&purchase.ID,
		&purchase.TotalPrice,
		&purchase.TotalTransfer,
		&purchase.SenderName,
		&purchase.SenderContactType,
		&purchase.SenderContactDetail,
		&purchase.PaidAt,
	)
	if err != nil {
		return model.Purchase{}, err
	}

	batch := &pgx.Batch{}
	for _, item := range arg.PurchasedItems {
		batch.Queue(insertPurchaseProductsQuery, purchase.ID, item.ProductID, item.Qty)
	}
	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return model.Purchase{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Purchase{}, err
	}

	return purchase, nil
}

const getPurchase = `-- name: GetPurchase :one
SELECT id, total_price, total_transfer, sender_name, sender_contact_type, sender_contact_detail, paid_at FROM purchases
WHERE id = $1
LIMIT 1
`

func (r *PurchaseRepository) GetPurchase(ctx context.Context, purchaseId int) (model.Purchase, error) {
	row := r.pool.QueryRow(ctx, getPurchase, purchaseId)
	var i model.Purchase
	err := row.Scan(
		&i.ID,
		&i.TotalPrice,
		&i.TotalTransfer,
		&i.SenderName,
		&i.SenderContactType,
		&i.SenderContactDetail,
		&i.PaidAt,
	)
	return i, err
}

const getPurchaseProductsByIdQuery = `-- name: GetPurchaseProducts :many
SELECT id, purchase_id, product_id, qty, created_at FROM pivot_purchase_products
WHERE purchase_id = $1
`

func (r *PurchaseRepository) GetPurchaseProductsById(ctx context.Context, purchaseId int) ([]model.PurchaseProduct, error) {
	rows, err := r.pool.Query(ctx, getPurchaseProductsByIdQuery,
		purchaseId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.PurchaseProduct
	for rows.Next() {
		var i model.PurchaseProduct
		if err := rows.Scan(
			&i.ID,
			&i.PurchaseID,
			&i.ProductID,
			&i.Qty,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePurchasePaidAtQuery = `-- name: UpdatePurchasePaidAt :exec
UPDATE purchases
SET paid_at = $1
WHERE id = $2
`

const updateProductQtyQuery = `-- name: UpdateProductQty :exec
UPDATE products SET qty = qty - $1 WHERE id = $2
`

type UpdatePurchaseParams struct {
	PurchaseID       int
	PurchaseProducts []model.PurchaseProduct
}

func (r *PurchaseRepository) UpdatePurchase(ctx context.Context, arg UpdatePurchaseParams) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updatePurchasePaidAtQuery, time.Now(), arg.PurchaseID)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, product := range arg.PurchaseProducts {
		batch.Queue(updateProductQtyQuery, product.Qty, product.ProductID)
	}
	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
