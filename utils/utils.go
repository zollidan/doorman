package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zollidan/doorman/database"
	"github.com/zollidan/doorman/models"
	"gorm.io/gorm"
)


func ExtractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", false
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", false
	}

	return parts[1], true
}

func IssueTokens(user *models.User, jwtSecret string, db *gorm.DB) (string, string, error) {
	accessTokenExpiration := time.Now().Add(15 * time.Minute)
	refreshTokenExpiration := time.Now().Add(7 * 24 * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"typ": "access",
		"exp": jwt.NewNumericDate(accessTokenExpiration),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"type": "refresh",
		"exp": jwt.NewNumericDate(refreshTokenExpiration),
	})

	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("error generating access token: %w", err)
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("error generating refresh token: %w", err)
	}

	refreshTokenModel := &models.RefreshToken{
		UserID: user.ID,
		Token: refreshTokenString,
		ExpiresAt: refreshTokenExpiration,
	}

	err = database.Create(db, refreshTokenModel)
	if err != nil {
		return "", "", fmt.Errorf("error saving refresh token: %w", err)
	}

	session := &models.Session{
		UserID: user.ID,
		Token: accessTokenString,
		ExpiresAt: accessTokenExpiration,
	}	

	err = database.Create(db, session)
	if err != nil {
		return "", "", fmt.Errorf("error saving session: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}