package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uuid.UUID `gorm:"primaryKey"`
	Username      string
	Email         string
	PasswordHash  string
	EmailVerified bool
	IsActive      bool
}

// type Client struct {
// }
