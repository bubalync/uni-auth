package jwtgen

import (
	"errors"
	"fmt"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	UserId uuid.UUID `json:"uid"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type TokenGenerator interface {
	GenerateAccessToken(user entity.User) (string, error)
	GenerateRefreshToken(user entity.User) (string, error)

	ParseAccessToken(tokenStr string) (*Claims, error)
	ParseRefreshToken(tokenStr string) (*Claims, error)
}

type JWTTokenGenerator struct {
	accessSignKey  string
	accessTokenTTL time.Duration

	refreshSignKey  string
	refreshTokenTTL time.Duration
}

func NewJwtTokenGenerator(accessSignKey, refreshSignKey string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		accessSignKey:   accessSignKey,
		refreshSignKey:  refreshSignKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (g *JWTTokenGenerator) GenerateAccessToken(user entity.User) (string, error) {
	return generateToken(user, g.accessSignKey, g.accessTokenTTL)
}

func (g *JWTTokenGenerator) GenerateRefreshToken(user entity.User) (string, error) {
	return generateToken(user, g.refreshSignKey, g.refreshTokenTTL)
}

func generateToken(user entity.User, secret string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserId: user.Id,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (g *JWTTokenGenerator) ParseAccessToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr, g.accessSignKey)
}

func (g *JWTTokenGenerator) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr, g.refreshSignKey)
}

func parseToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to cast token claims")
	}

	return claims, nil
}
