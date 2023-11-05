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

// Database interface for handling database operations
type Database interface {
	InitProductCollection() error
	InsertProduct(product models.Product) error
	FetchProductsByImageOrderAndIsNotCompressed(limit int) ([]models.Product, error)
	UpdateProduct(updatedProduct models.Product) error
}

// MongoDB struct implements the Database interface
type MongoDB struct {
	productCollection *mongo.Collection
}

// NewMongoDB creates a new MongoDB instance
func NewMongoDB(client *mongo.Client) *MongoDB {
	return &MongoDB{
		productCollection: client.Database("Zocket").Collection("products"),
	}
}

// InitProductCollection initializes the MongoDB collection for events.
func (m *MongoDB) InitProductCollection() error {
	// Get the MongoDB client
	client = GetClient()

	// Access the "events" collection in the "your-database-name" database
	productCollection := client.Database("Zocket").Collection("products")

	// Create an index on the "created_at" field for sorting and querying
	indexOptions := options.Index().SetBackground(true).SetSparse(true)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "is_compressed", Value: 1}},
		Options: indexOptions,
	}

	_, err := productCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return fmt.Errorf("error creating index: %v", err)
	}

	fmt.Println("Product collection initialized.")
	return nil
}

// InsertProduct inserts an event into the "events" collection.
func (m *MongoDB) InsertProduct(product models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := m.productCollection.InsertOne(context.TODO(), product)
	if err != nil {
		return fmt.Errorf("error inserting product to DB: %v", err)
	}
	return nil
}

func (m *MongoDB) FetchProductsByImageOrderAndIsNotCompressed(limit int) ([]models.Product, error) {
	// Define the options for sorting and limiting.
	opts := options.Find().SetSort(bson.D{{"images.created_at", 1}}).SetLimit(int64(limit))

	// Create a filter to include only products where isCompressed is false.
	filter := bson.D{{"is_compressed", false}}

	// Find the products with sorting, limiting, and the isCompressed filter.
	cursor, err := m.productCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error fetching products from DB: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, context.TODO())

	// Decode the products from the cursor.
	var products []models.Product
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, fmt.Errorf("error decoding products: %v", err)
	}

	return products, nil
}

func (m *MongoDB) UpdateProduct(updatedProduct models.Product) error {
	// Define the filter to identify the product by its ID.
	filter := bson.D{{"_id", updatedProduct.ID}}

	// Define the update operation using the BSON document.
	update := bson.D{
		{"$set", bson.D{{"is_compressed", updatedProduct.IsCompressed}, {"compressed_product_images", updatedProduct.CompressedProductImages}}},
	}

	// Set the options for the update operation.
	opts := options.Update().SetUpsert(false) // SetUpsert(false) means the update won't create a new document if the filter doesn't match.

	// Perform the update operation.
	result, err := m.productCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return fmt.Errorf("error updating product: %v", err)
	}

	// Check if the update affected any documents.
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no product updated, product not found or already compressed")
	}

	return nil
}
