package Tests

import (
	"encoding/json"
	"fmt"
	"log"
	initializer "main/Initializer"
	models "main/Models"
	routing "main/Routing"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCorrectSignUp(t *testing.T) {
	router := routing.SetupRouter()
	w := httptest.NewRecorder()
	var body struct {
		Email    string
		Password string
	}
	body.Email = "tr@gmail.com"
	body.Password = "12345678"
	fmt.Println("Password is in test ", body.Password)
	js, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(js)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	initializer.DB.Unscoped().Where("email = ?", body.Email).Delete(&models.Users{})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"disc\":\"Account created\"}", w.Body.String())

}
func TestInorrectSignUpBody(t *testing.T) {
	router := routing.SetupRouter()
	w := httptest.NewRecorder()
	var body struct {
		Email string
	}
	body.Email = "tr@gmail.com"
	js, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(js)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"wrong struct of req body\"}", w.Body.String())

}
func TestUserExsistsSignUp(t *testing.T) {
	router := routing.SetupRouter()
	w := httptest.NewRecorder()
	var body struct {
		Email    string
		Password string
	}
	var user models.Users
	res := initializer.DB.First(&user)
	if res.Error != nil {
		log.Fatal("Сначала добавьте пользовтеля в бд")
	}
	fmt.Println(user.Email)
	body.Email = user.Email
	body.Password = "12345678"
	js, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(js)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCorrectLogin(t *testing.T) {
	router := routing.SetupRouter()
	w := httptest.NewRecorder()
	var body struct {
		Email    string
		Password string
	}
	body.Email = "tr@gmail.com"
	body.Password = "12345678"
	var usr models.Users
	usr.Email = body.Email
	pas, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	usr.Password = string(pas)
	initializer.DB.Save(&usr)
	js, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(string(js)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	fmt.Println(w.Code)
	if w.Code != http.StatusOK {
		fmt.Println(w.Body)
		log.Fatal("Ошибка ")

	}
	var user models.Users
	user.Email = body.Email
	res := initializer.DB.Where("email = ?", body.Email).First(&user)
	initializer.DB.Unscoped().Where("email = ?", body.Email).Delete(&models.Users{})
	if res.Error != nil {
		log.Fatal("Ошибка получения данных")
	}
	err := godotenv.Load()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 10).Unix(),
	})
	tokenstr, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Fatal("Ошибка ")
	}
	act := w.Result().Cookies()
	var actck string
	for _, val := range act {
		if val.Name == "Authorization" {
			actck = val.Value
		}

	}

	fmt.Println("ID ", user.ID)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, tokenstr, actck)
}
func TestIncorrectLogin(t *testing.T) {
	router := routing.SetupRouter()
	w := httptest.NewRecorder()
	var body struct {
		Email    string
		Password string
	}
	body.Email = "tr@gmail.com"
	body.Password = "12345678"
	js, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(string(js)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"Invalid email or password\"}", w.Body.String())

}
