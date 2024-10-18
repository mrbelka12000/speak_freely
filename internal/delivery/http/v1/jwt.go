package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

type (
	myClaims struct {
		ID        int64 `json:"id"`
		IsRefresh bool  `json:"is_refresh,omitempty"`
		jwt.RegisteredClaims
	}
)

func (m myClaims) Validate() error {
	if m.ID == 0 {
		return errors.New("empty id")
	}

	return nil
}

func createTokens(id int64) (map[string]string, error) {

	claims := myClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	claims = myClaims{
		ID:        id,
		IsRefresh: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenString, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return map[string]string{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	}, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func refreshTokens(tokenString string) (map[string]string, error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	claims, ok := token.Claims.(*myClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token format")
	}

	return createTokens(claims.ID)
}

// Tokens refresh JWT tokens
func (h *Handler) Tokens(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("RefreshToken")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokens, err := refreshTokens(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.log.With("error", err).Error("invalid token")
		return
	}

	writeJson(w, tokens, http.StatusOK)
}
