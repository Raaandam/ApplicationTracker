package controllers

// remember to change the mock userID after authentication
import (
	"context"
	"net/http"
	"time"

	"github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/database"
	models "github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var applicationCollection *mongo.Collection = database.OpenCollection("application-analysis")

var validate = validator.New()

// mock user ID as of now
//it converts string to objectID(how it is actually stored in the database)
// _ ignores the error sent by this

func GetApplications() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := bson.ObjectIDFromHex("6641a1f0e13c4a1b2c3d4e5f")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var applications []models.Application

		cursor, err := applicationCollection.Find(ctx, bson.M{"user_id": userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})

			return
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &applications); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, applications)
	}
}

func GetApplication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		applicationID := c.Param("_id")

		if applicationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
			return
		}

		objID, err := bson.ObjectIDFromHex(applicationID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID could not be convereted"})
			return
		}

		var application models.Application

		err = applicationCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&application)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}

		c.JSON(http.StatusOK, application)
	}
}

func AddApllication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var application models.Application

		if err := c.ShouldBindJSON(&application); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(application); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation Failed", "details": err.Error()})
			return
		}

		result, err := applicationCollection.InsertOne(ctx, application)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add application"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
