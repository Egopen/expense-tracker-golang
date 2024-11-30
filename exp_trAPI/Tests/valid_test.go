package Tests

import (
	"log"
	initializer "main/Initializer"
	models "main/Models"
	routing "main/Routing"
	valid "main/Valid"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func generateToken(id uint, secretKey string, expirationTime time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(expirationTime).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}

func TestCheckValidToken(t *testing.T) {
	godotenv.Load()
	gin.SetMode(gin.TestMode)
	router := routing.SetupRouter()
	router.GET("/protected", valid.CheckToken, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	user := models.Users{Email: "effdadsa", Password: "sadsafsa"}
	initializer.DB.Save(&user)
	defer initializer.DB.Unscoped().Delete(&user)
	tokenString := generateToken(user.ID, os.Getenv("SECRET_KEY"), time.Hour)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: tokenString})
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"success"}`, w.Body.String())
}
func TestCheckExpiredToken(t *testing.T) {
	godotenv.Load()
	gin.SetMode(gin.TestMode)
	router := routing.SetupRouter()
	router.GET("/protected", valid.CheckToken, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	tokenString := generateToken(1, os.Getenv("SECRET_KEY"), -time.Hour*2)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: tokenString,
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 Unauthorized for expired token")
}

func TestCheckNoUserToken(t *testing.T) {
	godotenv.Load()
	gin.SetMode(gin.TestMode)
	router := routing.SetupRouter()
	router.GET("/protected", valid.CheckToken, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	tokenString := generateToken(0, os.Getenv("SECRET_KEY"), time.Hour)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: tokenString,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 Unauthorized for expired token")
}
