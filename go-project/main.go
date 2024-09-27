package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kr/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://lahiruheshan454:eiq7iCmSR9ACIUgq@cluster0.ekc4f.mongodb.net/data?retryWrites=true&w=majority")

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
		{"name", "lahiru heshan"},
		{"email", "lahiruheshan454@gmail.com"},
	}
	insertResult, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	showAllDocuments(collection)

	updateResult, err := collection.UpdateOne(context.TODO(), bson.D{{"name", "K.B.L.L Heshan"}}, bson.D{{"$set", bson.D{{"email", "example.com"}}}})
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
	}
	fmt.Printf("Updated %v document(s)\n", updateResult.ModifiedCount)

	searchDocument(collection, "K.B.L.L Heshan")

	deleteResult, err := collection.DeleteOne(context.TODO(), bson.D{{"name", "K.B.L.L Heshan"}})
	if err != nil {
		log.Fatalf("Failed to delete document: %v", err)
	}
	fmt.Printf("Deleted %v document(s)\n", deleteResult.DeletedCount)

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
		pretty.Println(result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Error occurred during cursor iteration: %v", err)
	}
}

func searchDocument(collection *mongo.Collection, name string) {
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&result)
	if err != nil {
		log.Printf("Failed to find document with name %s: %v\n", name, err)
		return
	}
	fmt.Println("Found document:")
	pretty.Println(result)
}
