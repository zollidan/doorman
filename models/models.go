package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username      string         `gorm:"uniqueIndex;not null"`
	Email         string         `gorm:"uniqueIndex;not null"`
	PasswordHash  string         `gorm:"not null"`
	EmailVerified bool           `gorm:"default:false"`
	IsActive      bool           `gorm:"default:true"`
	LastLoginAt   *time.Time     `gorm:"index"`
	Sessions      []Session      `gorm:"foreignKey:UserID"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`
}

type Session struct {
	gorm.Model
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	User      User       `gorm:"foreignKey:UserID"`
	Token     string     `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"not null;index"`
	IPAddress string     `gorm:"size:45"`
	UserAgent string     `gorm:"size:255"`
	IsActive  bool       `gorm:"default:true;index"`
	RevokedAt *time.Time `gorm:"index"`
}

type RefreshToken struct {
	gorm.Model
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	User      User       `gorm:"foreignKey:UserID"`
	Token     string     `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"not null;index"`
	IsUsed    bool       `gorm:"default:false;index"`
	RevokedAt *time.Time `gorm:"index"`
}

type LoginAttempt struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email      string    `gorm:"index;not null"`
	IPAddress  string    `gorm:"size:45;index"`
	Successful bool      `gorm:"index"`
	FailReason string    `gorm:"size:255"`
}