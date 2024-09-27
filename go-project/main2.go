package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Order struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Customer string             `bson:"customer"`
	Item     string             `bson:"item"`
	Quantity int                `bson:"quantity"`
	Price    float64            `bson:"price"`
}

func runOrderService() {

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
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("data").Collection("orders")

	order := Order{
		Customer: "Alice",
		Item:     "Pizza",
		Quantity: 2,
		Price:    15.50,
	}

	insertResult, err := collection.InsertOne(context.TODO(), order)
	if err != nil {
		log.Fatalf("Failed to insert order: %v", err)
	}
	fmt.Printf("Inserted order with ID: %v\n", insertResult.InsertedID)

	showAllOrders(collection)

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
}

func showAllOrders(collection *mongo.Collection) {
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to find orders: %v", err)
	}
	defer cursor.Close(context.TODO())

	fmt.Println("Orders in the collection:")
	for cursor.Next(context.TODO()) {
		var order Order
		err := cursor.Decode(&order)
		if err != nil {
			log.Fatalf("Failed to decode order: %v", err)
		}
		fmt.Printf("Order ID: %v, Customer: %s, Item: %s, Quantity: %d, Price: %.2f\n",
			order.ID, order.Customer, order.Item, order.Quantity, order.Price)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Error during cursor iteration: %v", err)
	}
}
