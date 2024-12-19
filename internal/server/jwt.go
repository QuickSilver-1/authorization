package server

import (
	"crypto/rand"
	"fmt"
	"time"

	e "auth/internal/errors"
	"auth/internal/logger"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Id int    `json:"id"`
	Ip string `json:"ip"`
	jwt.StandardClaims
}

func NewToken(id int, ip string) *Token {
	return &Token{
		Id: id,
		Ip: ip,
	}
}

// MakeToken создает JWT токен для данного email и id
func (t *Token) MakeJWT(key string) (string, error) {
	// Устанавливаем время истечения токена
	logger.Log.Debug("creating token")
	expires := time.Now().Add(15 * time.Minute)
	claims := &Token{
		Id: t.Id,
		Ip: t.Ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenS, err := token.SignedString([]byte(key))

	if err != nil {
		logger.Log.Error(fmt.Sprintf("generation token error: %v", err))
		return "", &e.JWTError{
			Err: fmt.Sprintf("Generation token error: %v", err),
		}
	}

	logger.Log.Debug("creating successful")
	return tokenS, nil
}

// DecodeJWT декодирует JWT токен и возвращает Token
func (t *Token) DecodeJWT(tokenStr string, key string) (*Token, error) {
	logger.Log.Debug("decoding start")
	claims := &Token{}
	// Парсинг токена с claims
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, &e.JWTError{
				Err: "Invalid form token",
			}
		}

		logger.Log.Error(fmt.Sprintf("Generation token error: %v", err))
		return nil, &e.JWTError{
			Err: fmt.Sprintf("Decode token error: %v", err),
		}
	}

	if !token.Valid {
		return nil, &e.JWTError{
			Err: "Invalid token",
		}
	}

	logger.Log.Debug("decoding success")
	return claims, nil
}

func (t *Token) CreateRefresh() ([]byte, error) {
	s := make([]byte, 64)
	_, err := rand.Read(s)

	if err != nil {
		return nil, &e.CryptError{
			Err: fmt.Sprintf("Refresh token generation error: %s", err),
		}
	}

	return s, nil
}
