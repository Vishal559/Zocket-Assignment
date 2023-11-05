# Zocket-Assignment

## Introduction 
The Aim of the assignment to build a product management system where authenticated user can create product. The system
is capable to process those product without blocking the user operation. Whiling processing the product, system compresses 
the product images and update the product schema with compress images. 

## Proposed Solution 
The solution provides a API to users to create products. Only authenticated user can create the products by providing 
product details & product images. Creating product with images is completely optional. There is a cron job which keep
running in the background every 1 minute. The cron job fetches 10 products which is not processed yet and process those
by compressing their images. 

Notes - 
A additional parameter added to the product schema. Modified Schema
Product Id considered as string rather than int.
`
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
`
## Image Compression
The product images get stored in the folder /uploads. The compressed images will get stored in compressedImageUploads 
folder. 

Notes - 
Only .jpg images supported by the image compression logic.

## Testing 
As only authenticated user can create the images. testUser1 is already onboarded to the database user_id - "1234".

## Unit Testing 
Unit Tests only added for product controllers. 
User Controller & Dao Method Unit Tests yet to add. 

## TESTING VIA POSTMAN
![request-response](https://drive.google.com/file/d/131-Wrh8xR1MgOpGPBvbhvOjJK-9WGohv/view?usp=sharing)
![request-headers](https://drive.google.com/file/d/1CwiHsIfrmkc2bSVn7FFagwBWePpw7CMl/view?usp=sharing)



