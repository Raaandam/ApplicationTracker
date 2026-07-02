package main

import (
	"fmt"

	controller "github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/applications", controller.GetApplications())
	router.GET("application/:_id", controller.GetApplication())
	router.POST("/:User_id/addapplication", controller.AddApllication())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
