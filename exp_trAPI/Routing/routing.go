package routing

import (
	controllers "main/Controllers"
	initializer "main/Initializer"
	valid "main/Valid"

	"github.com/gin-gonic/gin"
)

func Init() {
	initializer.LoadEnv()
	initializer.InitDB()
	initializer.Sync()
}
func SetupRouter() *gin.Engine {
	Init()
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/check", valid.CheckToken, controllers.Check)
	r.GET("/exp/get-all", valid.CheckToken, controllers.GetALlExp)
	r.POST("/exp/add", valid.CheckToken, controllers.AddExp)
	r.PUT("/exp/update", valid.CheckToken, controllers.UpdateExp)
	r.DELETE("/exp/delete", valid.CheckToken, controllers.DeleteExp)
	r.GET("/exp/get-last-month", valid.CheckToken, controllers.GetForLastMonth)
	r.GET("/exp/get-last-three-month", valid.CheckToken, controllers.GetForLastThreeMonth)
	r.GET("/exp/get-from-period", valid.CheckToken, controllers.GetForCustomPeriod)
	r.GET("/cat/get-all", valid.CheckToken, controllers.GetAllCat)
	return r
}
