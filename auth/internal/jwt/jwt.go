package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type UserInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type TokenManager struct {
	refreshSecretKey []byte
	accessSecretKey  []byte
	accessExpiry     time.Duration
	refreshExpiry    time.Duration
}

func NewTokenManager(refreshSecret, accessSecret string, refreshExpiry, accessExpiry time.Duration) *TokenManager {
	return &TokenManager{
		refreshSecretKey: []byte(refreshSecret),
		accessSecretKey:  []byte(accessSecret),
		refreshExpiry:    refreshExpiry,
		accessExpiry:     accessExpiry,
	}
}

func (tm *TokenManager) GenerateRefreshToken(info UserInfo) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.refreshExpiry)),
		},
		UserID:   info.UserID,
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.refreshSecretKey)
}

func (tm *TokenManager) GenerateAccessToken(info UserInfo) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessExpiry)),
		},
		UserID:   info.UserID,
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.accessSecretKey)
}

func (tm *TokenManager) VerifyRefreshToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}
			return tm.refreshSecretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid refresh token: %s", err.Error())
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

func (tm *TokenManager) VerifyAccessToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}
			return tm.accessSecretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid access token: %s", err.Error())
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid access token claims")
	}

	return claims, nil
}
