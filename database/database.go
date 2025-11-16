package database

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/zollidan/doorman/config"
	"github.com/zollidan/doorman/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Model interface {
	models.User | models.RefreshToken | models.Session | models.LoginAttempt
}

func Init(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_HOST, cfg.POSTGRES_PORT, cfg.POSTGRES_DB)
	db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Printf("Failed to connect to PostgreSQL: %v", err)
		log.Println("Falling back to SQLite database...")

		db, err = gorm.Open(sqlite.Open("doorman.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect to any database")
		}
		log.Println("Successfully connected to SQLite database")
	} else {
		log.Println("Successfully connected to PostgreSQL database")
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

func GetBy[T Model](db *gorm.DB, searchField string, value any) (*T, error) {
	result, err := gorm.G[T](db).Where(searchField+" = ?", value).First(context.Background())
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
