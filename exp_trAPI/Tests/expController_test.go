package Tests

import (
	"encoding/json"
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
)

func TestGetAll(t *testing.T) {
	godotenv.Load()
	r := routing.SetupRouter()
	var usr models.Users
	usr.Email = "tr@gmail.com"
	usr.Password = "12345678"
	initializer.DB.Save(&usr)

	var cat models.Categories
	cat.Name = "Кредит"
	initializer.DB.Save(&cat)

	var exp1, exp2 models.Expenses
	exp1 = models.Expenses{
		Date:       time.Now().Add(-(time.Hour * 100)),
		CategoryID: int(cat.ID),
		UserID:     int(usr.ID),
		Sum:        2000,
	}
	exp2 = models.Expenses{
		Date:       time.Now().Add(-(time.Hour * 100)),
		CategoryID: int(cat.ID),
		UserID:     int(usr.ID),
		Sum:        2000,
	}
	initializer.DB.Save(&exp1)
	initializer.DB.Save(&exp2)

	defer initializer.DB.Unscoped().Delete(&exp1)
	defer initializer.DB.Unscoped().Delete(&exp2)
	defer initializer.DB.Unscoped().Delete(&usr)
	defer initializer.DB.Unscoped().Delete(&cat)

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		t.Fatal("SECRET_KEY not set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.ID,
		"exp": time.Now().Add(time.Hour * 10).Unix(),
	})
	tokenstr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/exp/get-all", nil)

	req.AddCookie(&http.Cookie{
		Name:     "Authorization",
		Value:    tokenstr,
		MaxAge:   3600 * 10,
		HttpOnly: true,
	})
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
func TestAdd(t *testing.T) {
	godotenv.Load()
	r := routing.SetupRouter()
	var usr models.Users
	usr.Email = "tr@gmail.com"
	usr.Password = "12345678"
	initializer.DB.Save(&usr)

	var cat models.Categories
	cat.Name = "Кредит"
	initializer.DB.Save(&cat)

	var jsexp struct {
		Date        time.Time `json:"Date"`
		Sum         int       `json:"Sum"`
		Discription string    `json:"Discription"`
		CategoryID  int       `json:"CategoryID"`
	}
	jsexp.Date = time.Now().Add(-(time.Hour * 100))
	jsexp.CategoryID = int(cat.ID)
	jsexp.Sum = 2000
	defer initializer.DB.Unscoped().Delete(&usr)
	defer initializer.DB.Unscoped().Delete(&cat)

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		t.Fatal("SECRET_KEY not set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.ID,
		"exp": time.Now().Add(time.Hour * 10).Unix(),
	})
	tokenstr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	js, _ := json.Marshal(jsexp)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/exp/add", strings.NewReader(string(js)))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:     "Authorization",
		Value:    tokenstr,
		MaxAge:   3600 * 10,
		HttpOnly: true,
	})
	r.ServeHTTP(w, req)
	var jsexpres struct {
		Exp models.Expenses `json:"exp"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &jsexpres)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	res := initializer.DB.Unscoped().Where("id= ?", jsexpres.Exp.ID).Delete(&models.Expenses{})
	if res.Error != nil {
		t.Fatalf("Error: %v", err)
	}
	assert.Equal(t, http.StatusOK, w.Code)

}
