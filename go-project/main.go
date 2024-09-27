package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatalf("MONGODB_URI not set in .env file")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("data").Collection("data")

	document := bson.D{
		{"name", "John Doe"},
		{"email", "john.doe@example.com"},
	}
	insertResult, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("********************************************************* \n")
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
	fmt.Printf("********************************************************* \n")
	showAllDocuments(collection)

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from the database: %v", err)
		}
	}()
}

func showAllDocuments(collection *mongo.Collection) {
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to find documents: %v", err)
	}
	defer cursor.Close(context.TODO())

	fmt.Println("Documents in the collection:")
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatalf("Failed to decode document: %v", err)
		}
		fmt.Printf("********************************************************* \n")
		fmt.Println(result)
		fmt.Printf("********************************************************* \n")
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Error occurred during cursor iteration: %v", err)
	}
}
