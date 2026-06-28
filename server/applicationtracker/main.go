package main

import (
	"fmt"

	controller "github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/applications", controller.GetApplications())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
