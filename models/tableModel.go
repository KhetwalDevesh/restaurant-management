package models

import "time"

type Table struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	NumberOfGuests uint      `json:"numberOfGuests"`
	TableNumber    uint      `json:"tableNumber"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
