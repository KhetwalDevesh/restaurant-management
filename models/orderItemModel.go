package models

import "time"

type Quantity string

const (
	S Quantity = "S"
	M Quantity = "M"
	L Quantity = "L"
)

type OrderItem struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Quantity  Quantity  `json:"quantity" gorm:"not null" validate:"oneof=S M L"`
	UnitPrice float64   `json:"unitPrice" gorm:"not null"`
	FoodId    uint      `gorm:"not null" json:"foodId"`
	OrderId   uint      `gorm:"not null" json:"orderId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
