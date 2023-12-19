package models

import "time"

type Order struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	OrderDate time.Time `gorm:"not null" json:"orderDate"`
	TableId   uint32    `gorm:"not null" json:"tableId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
