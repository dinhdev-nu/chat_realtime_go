package jwt

import (
	"fmt"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	jwt.RegisteredClaims
}

func CreateToken(uuidToken string) (string, error) {

	ttl := time.Now().Add(time.Duration(global.Config.Jwt.JwtExpireTime) * time.Hour)
	payload := &Payload{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authgo",
			Subject:   uuidToken,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(ttl),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.Jwt.JwtSecret))
}

func VerifyToken(token string) (*jwt.RegisteredClaims, error) {
	// parse token
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(global.Config.Jwt.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
