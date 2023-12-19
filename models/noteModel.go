package models

import "time"

type Note struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Text      string    `json:"text"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
