package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zollidan/doorman/config"
	"github.com/zollidan/doorman/database"
	"github.com/zollidan/doorman/models"
	"github.com/zollidan/doorman/schemas"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handlers struct {
	cfg *config.Config
	db *gorm.DB
}

func New(cfg *config.Config, db *gorm.DB) *Handlers {
	return &Handlers{cfg: cfg, db: db}
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

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req *schemas.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := database.GetByEmail[models.User](h.db, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			database.Create(h.db, &models.LoginAttempt{
				Email: req.Email,
				Successful: false,
				FailReason: "Invalid email",
			})
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("Error fetching user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(req.Password))
	if err != nil {
		database.Create(h.db, &models.LoginAttempt{
			Email: req.Email,
			Successful: false,
			FailReason: "Invalid password",
		})
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	accessTokenExpiration := time.Now().Add(15 * time.Minute)
	refreshTokenExpiration := time.Now().Add(7 * 24 * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": result.ID,
		"email":   result.Email,
		"typ": "access",
		"exp": jwt.NewNumericDate(accessTokenExpiration),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": result.ID,
		"email":   result.Email,
		"typ": "refresh",
		"exp": jwt.NewNumericDate(refreshTokenExpiration),
	})

	accessTokenString, err := accessToken.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating refresh token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	refreshTokenModel := &models.RefreshToken{
		UserID: result.ID,
		Token: refreshTokenString,
		ExpiresAt: refreshTokenExpiration,
	}

	err = database.Create(h.db, refreshTokenModel)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving refresh token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	session := &models.Session{
		UserID: result.ID,
		Token: accessTokenString,
		ExpiresAt: accessTokenExpiration,
	}	

	err = database.Create(h.db, session)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving session: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	resp := &schemas.TokenResponse{
		AccessToken: accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType: "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}