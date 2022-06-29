package main

import (
	"assignment2/controllers"
	"assignment2/database"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"assignment2/docs"
)

func init() {
	database.ConnectDB()
}

func main() {
	r := gin.Default()

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8888"
	docs.SwaggerInfo.Schemes = []string{"http"}

	ctrl := controllers.NewCarsController(database.GetDB_Order(), database.GetDB_Item())

	r.POST("/orders", ctrl.CreateOrders)
	r.GET("/orders", ctrl.GetAllOrders)
	r.PUT("/orders/:orderId", ctrl.UpdateOrders)
	r.DELETE("/orders/:orderId", ctrl.DeleteOrders)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8888")
}
