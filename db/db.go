package db

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	var uri = os.Getenv("MONGODB_URI")
	var name = os.Getenv("DB_NAME")
	clientOptions := options.Client().ApplyURI(uri + "/" + name)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetCollection(collectionName string) *mongo.Collection {
	client, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())
	erro := godotenv.Load(".env")
	if erro != nil {
		panic(erro)
	}
	var name = os.Getenv("DB_NAME")
	collection := client.Database(name).Collection(collectionName)
	return collection
}
