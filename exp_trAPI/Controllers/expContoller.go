package controllers

import (
	features "main/Features"
	initializer "main/Initializer"
	models "main/Models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetALlExp(c *gin.Context) {
	id, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var expenses []models.Expenses
	res := initializer.DB.Where("user_id = ?", id).Find(&expenses)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"expenses": expenses,
	})
}
func AddExp(c *gin.Context) {
	var jsexp struct {
		Date        time.Time `json:"Date" binding:"required"`
		Sum         int       `json:"Sum" binding:"required"`
		Discription string
		CategoryID  int `json:"CategoryID" binding:"required"`
	}
	if c.Bind(&jsexp) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong body of request",
		})
		return
	}
	id, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var exp = models.Expenses{Date: jsexp.Date, Discription: &jsexp.Discription, Sum: jsexp.Sum, CategoryID: jsexp.CategoryID, UserID: id}
	res := initializer.DB.Save(&exp)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to save data",
		})
		return
	}
	var jsexpres struct {
		ID          int       `json:"ID"`
		Date        time.Time `json:"Date"`
		Sum         int       `json:"Sum"`
		Discription string
		CategoryID  int `json:"CategoryID"`
	}
	jsexpres.ID = int(exp.ID)
	jsexpres.Date = exp.Date
	jsexpres.Sum = exp.Sum
	jsexpres.Discription = *exp.Discription
	jsexpres.CategoryID = exp.CategoryID
	c.JSON(http.StatusOK, gin.H{
		"exp": jsexpres,
	})

}
func UpdateExp(c *gin.Context) {
	var jsexp struct {
		ID          int       `json:"Id" binding:"required"`
		Date        time.Time `json:"Date" binding:"required"`
		Sum         int       `json:"Sum" binding:"required"`
		Discription *string
		CategoryID  int `json:"CategoryId" binding:"required"`
	}
	if c.Bind(&jsexp) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong body of request",
		})
		return
	}
	id, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var exp models.Expenses
	res := initializer.DB.First(&exp, "id = ? AND user_id = ?", strconv.Itoa(jsexp.ID), strconv.Itoa(id))
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find expense",
		})
		return
	}
	exp.Date = jsexp.Date
	exp.CategoryID = jsexp.CategoryID
	exp.Discription = jsexp.Discription
	exp.Sum = jsexp.Sum
	res = initializer.DB.Save(&exp)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to save expense",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"exp": exp,
	})
}
func DeleteExp(c *gin.Context) {
	var js struct {
		ID int `json:"Id" binding:"required"`
	}
	if c.Bind(&js) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong body of request",
		})
		return
	}
	id, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var exp models.Expenses
	exp.ID = uint(js.ID)
	exp.UserID = id
	res := initializer.DB.Unscoped().Delete(&exp)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete expense",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Expense": exp,
	})
}
func GetForLastMonth(c *gin.Context) {
	id, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var exps []models.Expenses
	res := initializer.DB.Where("date > ? AND user_id = ?", time.Now().UTC().AddDate(0, -1, 0), strconv.Itoa(id)).Find(&exps)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find expenses",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Expenses": exps,
	})
}
func GetForLastThreeMonth(c *gin.Context) {
	id, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var exps []models.Expenses
	res := initializer.DB.Where("date > ? AND user_id = ?", time.Now().UTC().AddDate(0, -3, 0), strconv.Itoa(id)).Find(&exps)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find expenses",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Expenses": exps,
	})
}
func GetForCustomPeriod(c *gin.Context) {
	From := c.Query("from")
	To := c.Query("to")
	if From == "" || To == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong path of request",
		})
		return
	}
	var js struct {
		From time.Time
		To   time.Time
	}
	var err error
	js.From, err = time.Parse("2006-01-02", From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong path of request",
		})
		return
	}
	js.To, err = time.Parse("2006-01-02", To)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong path of request",
		})
		return
	}
	if js.From.Unix() > js.To.Unix() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong date period",
		})
		return
	}
	id, err := features.GetIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var exps []models.Expenses
	res := initializer.DB.Where("date > ? AND date < ? AND user_id = ?", js.From.UTC(), js.To.UTC(), strconv.Itoa(id)).Find(&exps)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find expenses",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Expenses": exps,
	})
}
