package controllers

import (
	features "main/Features"
	initializer "main/Initializer"
	models "main/Models"
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"Email" binding:"required"`
		Password string `json:"Password" binding:"required"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong struct of req body",
		})
		return
	}
	email, err := mail.ParseAddress(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong format of email",
		})
		return
	}
	var person models.Users
	res := initializer.DB.Where(" email = ?", email.Address).First(&person)
	if res.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Here is account with the same email",
		})
		return
	}
	if len([]rune(body.Password)) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must contains more than 7 symbols",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cant make hash of password",
		})
		return
	}
	user := models.Users{Email: body.Email, Password: string(hash)}
	result := initializer.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cant save data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"disc": "Account created",
	})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"Email" binding:"required"`
		Password string `json:"Password" binding:"required"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong struct of req body",
		})
		return
	}
	email, err := mail.ParseAddress(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	var user models.Users
	res := initializer.DB.First(&user, "email=?", email.Address)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 10).Unix(),
	})
	tokenstr, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenstr, 3600*10, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Check(c *gin.Context) {
	user, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no such user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": user,
	})

}
