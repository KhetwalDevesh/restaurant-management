package models

import "time"

type Food struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	Image     string    `gorm:"" json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	MenuId    uint32    `json:"menuId" gorm:"not null"`
	// Association
	Menu Menu `gorm:"foreignKey:MenuId"`
}
