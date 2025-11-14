package database

import (
	"context"

	"github.com/zollidan/doorman/config"
	"github.com/zollidan/doorman/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Model interface {
	models.User | models.RefreshToken | models.Session | models.LoginAttempt
}

func Init(cfg *config.Config) *gorm.DB {

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.RefreshToken{},
		&models.LoginAttempt{},
	)

	return db
}


func Create[T Model](db *gorm.DB, item *T) error {
	err := gorm.G[T](db).Create(context.Background(), item)
	if err != nil {
		return err
	}
	return nil
}

func GetByEmail[T Model](db *gorm.DB, email string) (*T, error) {
	result, err := gorm.G[T](db).Where("email = ?", email).First(context.Background())
	if err != nil {
		return nil, err
	}
	return &result, nil
}