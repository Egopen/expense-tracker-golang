package Tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	features "main/Features"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetIdFromCookie(t *testing.T) {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла:", err)
	}

	// Проверяем наличие SECRET_KEY
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY не установлен в .env файле")
	}
	fmt.Println("Значение SECRET_KEY:", secretKey)

	// Создаем новый тестовый запрос
	req := httptest.NewRequest("GET", "/", nil)

	// Создаем новый тестовый контекст
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Генерируем токен
	want := 10
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": want,
	})

	// Подписываем токен
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Ошибка при создании токена: %v", err)
	}
	fmt.Println("Сгенерированный токен:", tokenStr)

	// Устанавливаем куку в запрос
	c.Request.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: tokenStr,
		Path:  "/",
	})

	// Проверяем, установлена ли кука
	cookieValue, err := c.Cookie("Authorization")
	if err != nil {
		t.Fatalf("Ошибка при получении куки: %v", err)
	}
	fmt.Printf("Кука установлена: Authorization = %s\n", cookieValue)

	// Вызов функции и проверка результата
	act, err := features.GetIdFromCookie(c)
	if err != nil {
		t.Fatalf("Ошибка при получении ID из куки: %v", err)
	}
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, want, act)
}
func TestGetIdFromCookieNoCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	want := 0
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	act, err := features.GetIdFromCookie(c)
	assert.NotErrorIs(t, err, nil)
	assert.Equal(t, want, act)

}

func TestGetIdFromCookieNoSub(t *testing.T) {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла:", err)
	}

	// Проверяем наличие SECRET_KEY
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY не установлен в .env файле")
	}
	fmt.Println("Значение SECRET_KEY:", secretKey)

	// Создаем новый тестовый запрос
	req := httptest.NewRequest("GET", "/", nil)

	// Создаем новый тестовый контекст
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Генерируем токен
	want := 0
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})

	// Подписываем токен
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Ошибка при создании токена: %v", err)
	}
	fmt.Println("Сгенерированный токен:", tokenStr)

	// Устанавливаем куку в запрос
	c.Request.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: tokenStr,
		Path:  "/",
	})

	// Проверяем, установлена ли кука
	cookieValue, err := c.Cookie("Authorization")
	if err != nil {
		t.Fatalf("Ошибка при получении куки: %v", err)
	}
	fmt.Printf("Кука установлена: Authorization = %s\n", cookieValue)

	act, err := features.GetIdFromCookie(c)
	assert.NotErrorIs(t, err, nil)
	assert.Equal(t, want, act)

}
