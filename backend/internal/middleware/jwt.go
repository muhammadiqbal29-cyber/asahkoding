package middleware

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var getSecretKey = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("default_fallback_secret_do_not_use_in_prod")
	}
	return []byte(secret)
}

// Claims adalah data apa saja yang ingin diselipkan ke dalam Token JWT
type Claims struct {
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken membuat token JWT baru saat User berhasil Login
func GenerateToken(userID uint, role string, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku 24 jam
	
	claims := &Claims{
		UserID:   userID,
		Role:     role,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Stempel token dengan Kunci Rahasia
	return token.SignedString(getSecretKey())
}

// ValidateToken memeriksa apakah token asli buatan server kita dan belum kedaluwarsa
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
