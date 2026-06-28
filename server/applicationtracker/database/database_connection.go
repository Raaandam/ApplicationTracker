package database

// remember to change the nil condition in DBInstance
import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env")
	}

	MongoDB := os.Getenv("MONGODB_URI")
	if MongoDB == "" {
		log.Fatal("MONGODB_URI not set")
	}

	fmt.Println("MONGODB_URI: ", MongoDB)

	clientOptions := options.Client().ApplyURI(MongoDB)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil // keeping it nil for now
	}

	return client
}

var Client *mongo.Client = DBInstance()

// will give me a collection like application etc.
func OpenCollection(CollectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	if databaseName == "" {
		log.Fatal("DATABASE_NAME env var is not set")
	}

	fmt.Println("DATABASE_NAME: ", databaseName)

	collection := Client.Database(databaseName).Collection(CollectionName)

	return collection
}
