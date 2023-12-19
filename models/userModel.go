package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// User represents the user model.
type User struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"not null;unique" json:"email" validate:"email"`
	Password     string    `gorm:"not null" json:"-"`
	IsAdmin      bool      `gorm:"default:false" json:"isAdmin"`
	Token        string    `gorm:"-" json:"token,omitempty"`
	RefreshToken string    `gorm:"-" json:"refreshToken,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// BeforeCreate is a GORM hook that sets the CreatedAt and UpdatedAt fields.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

// BeforeUpdate is a GORM hook that sets the UpdatedAt field.
func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}
