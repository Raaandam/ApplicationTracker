package utils

//this is just logi/utility layer. It's only job is to generate tokens and update database
import (
	"context"
	"os"
	"time"

	"github.com/Raaandam/ApplicationTracker/Server/ApplicationTrackerServer/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// jwt installation needed
type SignedDetails struct {
	Email  string
	Name   string
	UserID string

	jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")
var userCollection *mongo.Collection = database.OpenCollection("user")

func GenerateAllTokens(email, name, userID string) (string, string, error) {
	claims := &SignedDetails{
		Email:  email,
		Name:   name,
		UserID: userID,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ApplicationTracker",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	//it should technically have only userID as it's only job is to get a new access token
	refreshClaims := &SignedDetails{
		Email:  email,
		Name:   name,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ApplicationTracker",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func UpdateAllToken(userID, token, refreshToken string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": objID}, updateData)
	if err != nil {
		return err
	}
	return nil
}
