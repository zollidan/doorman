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
	"github.com/zollidan/doorman/utils"
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

	result, err := database.GetBy[models.User](h.db, "email", req.Email)
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

	accessTokenString, refreshTokenString, err := utils.IssueTokens(result, h.cfg.JWTSecret, h.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error issuing tokens: %s", err.Error()), http.StatusInternalServerError)
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

func (h *Handlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req *schemas.RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	tokenInDatabase, err := database.GetBy[models.RefreshToken](h.db, "token", req.RefreshToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("Error fetching refresh token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if time.Now().After(tokenInDatabase.ExpiresAt) {
		http.Error(w, "Refresh token has expired", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing token: %s", err.Error()), http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		http.Error(w, "Invalid token type", http.StatusUnauthorized)
		return
	}

	accessTokenString, refreshTokenString, err := utils.IssueTokens(&tokenInDatabase.User, h.cfg.JWTSecret, h.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error issuing tokens: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = database.Delete[models.RefreshToken](h.db, tokenInDatabase.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting refresh token: %s", err.Error()), http.StatusInternalServerError)
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