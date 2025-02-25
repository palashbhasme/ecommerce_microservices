package utils

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/models"
)

func ParseToken(tokenString, jwtSecret string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil

}
