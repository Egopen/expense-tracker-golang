package features

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetIdFromCookie(c *gin.Context) (int, error) {
	if c.Request == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return 0, fmt.Errorf("не удалось получить запрос")
	}
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return 0, fmt.Errorf("не удалось получить куки: %v", err)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга токена: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Something went wrong")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("ошибка: поле `sub` отсутствует или имеет неверный тип")
	}
	return int(sub), nil

}
