package controllers

import (
	initializer "main/Initializer"
	models "main/Models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCat(c *gin.Context) {

	var dbcat []models.Categories
	res := initializer.DB.Find(&dbcat)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find categories",
		})
		return
	}
	var cat = make([]struct {
		Id   int
		Name string
	}, len(dbcat))
	for i := 0; i < len(dbcat); i++ {
		cat[i].Id = int(dbcat[i].ID)
		cat[i].Name = dbcat[i].Name
	}
	c.JSON(http.StatusOK, gin.H{
		"Categories": cat,
	})

}
