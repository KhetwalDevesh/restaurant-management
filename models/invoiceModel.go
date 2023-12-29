package models

import (
	"time"
)

// PaymentMethod represents the payment method enum.
type PaymentMethod string

const (
	Card PaymentMethod = "card"
	Cash PaymentMethod = "cash"
)

type PaymentStatus string

const (
	Pending PaymentStatus = "pending"
	Paid    PaymentStatus = "paid"
)

type Invoice struct {
	ID             uint32        `gorm:"primary_key" json:"id"`
	OrderID        uint32        `gorm:"not null" json:"orderID"`
	PaymentMethod  PaymentMethod `gorm:"not null" json:"paymentMethod" validate:"oneof=card cash"`
	PaymentStatus  PaymentStatus `gorm:"not null" json:"paymentStatus" validate:"oneof=pending paid"`
	PaymentDueDate time.Time     `json:"paymentDueDate"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}
