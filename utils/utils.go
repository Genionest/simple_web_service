package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	// 密码hash
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

func GenerateJWT(username string) (string, error) {
	// 单用户验证
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		// 过期时间
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	// 加密
	signed_token, err := token.SignedString([]byte("secret"))
	return "Bearer " + signed_token, err
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	// 解密
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username claims is not a string")
		}
		return username, nil
	}
	return "", err
}
