package db

import (
	"Zocker-Assignment/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// UserDatabase interface for handling user database operations
type UserDatabase interface {
	InitUserCollection() error
	InsertUser(user models.User) error
	FindUser(userId string) (*models.User, error)
}

// MongoDBUser struct implements the UserDatabase interface
type MongoDBUser struct {
	userCollection *mongo.Collection
}

// NewMongoDBUser creates a new MongoDBUser instance
func NewMongoDBUser(client *mongo.Client) *MongoDBUser {
	return &MongoDBUser{
		userCollection: client.Database("Zocket").Collection("users"),
	}
}

// Implement the UserDatabase interface methods
func (m *MongoDBUser) InitUserCollection() error {
	// Get the MongoDB client
	client = GetClient()

	// Access the "events" collection in the "your-database-name" database
	m.userCollection = client.Database("Zocket").Collection("users")

	// Create an index on the "user_id" field for sorting and querying
	indexOptions := options.Index().SetBackground(true).SetSparse(true)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: indexOptions,
	}

	_, err := m.userCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return fmt.Errorf("error creating index: %v", err)
	}

	fmt.Println("User collection initialized..")
	return nil
}

func (m *MongoDBUser) InsertUser(user models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := m.userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("error inserting event: %v", err)
	}
	return nil
}

func (m *MongoDBUser) FindUser(userId string) (*models.User, error) {
	var user models.User
	// Create a filter to find the user by ID
	filter := bson.M{"user_id": userId}

	err := m.userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("no user exists with id: %v", err)
	}
	return &user, nil
}
