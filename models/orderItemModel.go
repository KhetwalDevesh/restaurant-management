package models

import "time"

type Quantity string

const (
	S Quantity = "S"
	M Quantity = "M"
	L Quantity = "L"
)

type OrderItem struct {
	ID        uint32  `gorm:"primary_key" json:"id"`
	Quantity  uint32  `json:"quantity" gorm:"not null"`
	UnitPrice float64 `json:"unitPrice" gorm:"not null"`
	FoodId    uint32  `gorm:"not null" json:"foodId"`
	OrderId   uint32  `gorm:"not null" json:"orderId"`
	// Associations
	Food      Food      `gorm:"foreignKey:FoodId;onDelete:CASCADE"`
	Order     Order     `gorm:"foreignKey:OrderId;onDelete:CASCADE"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
