package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in MongoDB.
type Product struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ProductID               string             `bson:"product_id,omitempty" json:"product_id"`
	ProductName             string             `bson:"product_name" json:"product_name"`
	ProductDescription      string             `bson:"product_description" json:"product_description"`
	ProductImages           []string           `bson:"product_images" json:"product_images"`
	ProductPrice            int                `bson:"product_price" json:"product_price"`
	CompressedProductImages []string           `bson:"compressed_product_images" json:"compressed_product_images"`
	IsCompressed            bool               `bson:"is_compressed" json:"is_compressed"`
	CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time          `bson:"updated_at" json:"updated_at"`
}
