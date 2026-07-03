package controllers

import (
	"github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

//database. calls the database package
var userCollection *mongo.Collection = database.OpenCollection("user")

func HashPassword(password string) (string, error){
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	return string(HashPassword)

}
