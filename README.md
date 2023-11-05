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

### Notes - 
A additional parameter added to the product schema.<br>
Product Id considered as string rather than int.<br>
<br>
Modified Schema
<br>
<br>
<br>
ID                      &nbsp;                           primitive.ObjectID <br>
ProductID               &nbsp;                           string            
ProductName             &nbsp;                           string            
ProductDescription      &nbsp;                           string             
ProductImages           &nbsp;                           []string           
ProductPrice            &nbsp;                           int                
CompressedProductImages &nbsp;                           []string           
IsCompressed            &nbsp;                           bool               
CreatedAt               &nbsp;                           time.Time          
UpdatedAt               &nbsp;                           time.Time          

## Image Compression
The product images get stored in the folder /uploads. The compressed images will get stored in compressedImageUploads 
folder. 

### Notes - 
Only .jpg images supported by the image compression logic.

## Testing 
As only authenticated user can create the images. testUser1 is already onboarded to the database user_id - "1234".

## Unit Testing 
Unit Tests only added for product controllers. 
User Controller & Dao Method Unit Tests yet to add. 

## TESTING VIA POSTMAN
! [request-response](https://drive.google.com/file/d/131-Wrh8xR1MgOpGPBvbhvOjJK-9WGohv/view)
! [request-headers](https://drive.google.com/file/d/1CwiHsIfrmkc2bSVn7FFagwBWePpw7CMl/view)



