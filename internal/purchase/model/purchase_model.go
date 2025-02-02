package model

import "time"

type Purchase struct {
	ID                  int
	TotalPrice          int
	TotalTransfer       int
	SenderName          string
	SenderContactType   string
	SenderContactDetail string
	PaidAt              *time.Time
}

type PurchaseProduct struct {
	ID         int
	PurchaseID int
	ProductID  int
	Qty        int
	CreatedAt  time.Time
}
