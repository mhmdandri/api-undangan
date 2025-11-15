package auth

import (
	"api-undangan/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateToken(userID string, role string)(string, error){
	claims := jwt.MapClaims{
		"sub": userID,
		"role": role,
		"iss": config.Cfg.JWTIssuer,
		"aud": config.Cfg.JWTAudience,
		"exp": time.Now().Add(config.Cfg.JWTExpiresIn).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}