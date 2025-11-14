package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zollidan/doorman/database"
	"github.com/zollidan/doorman/models"
	"github.com/zollidan/doorman/schemas"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handlers struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handlers {
	return &Handlers{db: db}
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req *schemas.RegisterUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error hashing password: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = database.Create(h.db, user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
			
}