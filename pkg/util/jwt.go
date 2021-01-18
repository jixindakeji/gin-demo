package util

import (
	"gin-demo/pkg/setting"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.GConfig.APP.JwtSecret)

//Claims for jwt
type Claims struct {
	UserName string `json:"username"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}

//GenerateToken generate token
func GenerateToken(username string, id uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * 7 * time.Hour)

	claims := Claims{
		username,
		id,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "waf-admin",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

//ParseToken parse token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
