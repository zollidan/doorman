package database

import (
	"context"

	"github.com/google/uuid"
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

func GetByToken[T Model](db *gorm.DB, token string) (*T, error) {
	result, err := gorm.G[T](db).Where("token = ?", token).Preload("User", nil).First(context.Background())
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func Delete[T Model](db *gorm.DB, id uuid.UUID) (int, error) {
	result, err := gorm.G[T](db).Where("id = ?", id).Delete(context.Background())
	if err != nil {
		return 0, err
	}
	
	return result, nil
}