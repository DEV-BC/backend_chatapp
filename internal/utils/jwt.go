package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func InitJWT(key string) {
	jwtKey = []byte(key)
}

type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Name     string `json:"name"`
	Platform string `json:"X-Platform"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId int64, name, platform string) (string, error) {
	//expiration time
	exp := time.Now().Add(2 * time.Hour)

	if platform != "web" && platform != "mobile" {
		return "", errors.New("invalid platform for token")
	}

	//have a proper platform and our exp
	claims := &CustomClaims{
		UserID:   userId,
		Name:     name,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   fmt.Sprint(userId),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey) //this makes a cryptic string
}

func VerifyJWT(tokenStr string) (int64, string, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(t *jwt.Token) (any, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtKey, nil
		},
	)
	if err != nil {
		return 0, "", "", fmt.Errorf("token parse error: %v", err)
	}
	if !token.Valid {
		return 0, "", "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, "", "", fmt.Errorf("token parse error: %v", err)
	}
	if claims.UserID == 0 || claims.Name == "" || (claims.Platform != "web" && claims.Platform != "mobile") {
		return 0, "", "", fmt.Errorf("missing or invalid user claims error: %v", err)
	}

	return claims.UserID, claims.Name, claims.Platform, nil
}
