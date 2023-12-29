package models

import "time"

type Menu struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Category  string    `json:"category"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
