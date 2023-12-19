package models

import "time"

type Food struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	Image     string    `gorm:"not null" json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	FoodId    uint      `json:"foodId"`
	MenuId    uint      `json:"menuId" gorm:"not null"`
}
