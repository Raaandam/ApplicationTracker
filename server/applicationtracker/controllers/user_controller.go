package controllers

//change secure = true
import (
	"context"
	"net/http"
	"time"

	"github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/database"
	model "github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/models"
	"github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// database. calls the database package
var userCollection *mongo.Collection = database.OpenCollection("user")

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Format"})
			return
		}

		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation Failed"})
			return
		}

		hashedPassword, err := HashPassword(user.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to Hash Password"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		//name and other things will be already filled by ShouldBindJSON
		user.ID = bson.NewObjectID()
		user.CreatedAt = time.Now()
		user.Password = hashedPassword
		user.UpdatedAt = time.Now()

		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin model.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var foundUser model.User

		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.Name, foundUser.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		err = utils.UpdateAllTokens(foundUser.ID.Hex(), token, refreshToken)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		//c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
		//path should be "/" so that it is allowed everywhere instead of lets say /login only

		c.SetCookie("access_token", token, 3600, "/", "localhost", false, true)

		c.SetCookie("refresh_token", refreshToken, 86400, "/", "localhost", false, true)

		//secure=false for local development (change to true in production with HTTPS)
		//will see httpOnly = true later on

		c.JSON(http.StatusOK, model.UserResponse{
			UserID: foundUser.ID.Hex(),
			Name:   foundUser.Name,
			Email:  foundUser.Email,
		})
	}
}
