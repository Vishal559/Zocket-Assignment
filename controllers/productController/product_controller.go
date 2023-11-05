package productController

import (
	"Zocker-Assignment/models"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// ProductDatabase interface for handling product database operations
type ProductDatabase interface {
	InitProductCollection() error
	InsertProduct(product models.Product) error
	// will Add other methods as needed
}
type UserDatabase interface {
	FindUser(userId string) (*models.User, error)
}

func CreateProduct(w http.ResponseWriter, r *http.Request, productDB ProductDatabase, userDB UserDatabase) {
	var product models.Product
	// Parse the form data, which may include image uploads
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form data: %v", err), http.StatusBadRequest)
		return
	}
	// Handle Validations
	if r.FormValue("user_id") == "" || r.FormValue("product_name") == "" || r.FormValue("product_description") == "" {
		http.Error(w, "userid, product name, product description & product price are required fields", http.StatusBadRequest)
		return
	}
	productPrice, err := strconv.ParseInt(r.FormValue("product_price"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing the product price: %v", err), http.StatusInternalServerError)
		return
	}
	if productPrice < 0 {
		http.Error(w, "product price cannot be negative", http.StatusBadRequest)
		return
	}
	_, err = userDB.FindUser(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, "only valid user can create products..try with user id (1234)", http.StatusBadRequest)
	}

	var productImages []string

	// Get the file headers for the specified key
	fileHeaders, ok := r.MultipartForm.File["images"]
	if !ok {
		fmt.Println("No files found for key: images")
	} else {
		// Handle image uploads
		for _, fileHeaders = range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				productImage, productImageErr := handleImage(fileHeader)
				if productImageErr != nil {
					http.Error(w, fmt.Sprintf("Error processing the images: %v", err), http.StatusInternalServerError)
					return
				}
				productImages = append(productImages, productImage)
			}
		}
	}

	product.ProductName = r.FormValue("product_name")
	product.ProductDescription = r.FormValue("product_description")
	product.ProductPrice = int(productPrice)
	product.ProductImages = productImages

	productId := uuid.New().String()
	product.ProductID = productId

	err = productDB.InsertProduct(product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting product into MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	// Return a success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		return
	}
}

// Define a function to handle each image part
func handleImage(fileHeader *multipart.FileHeader) (string, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	// Save the uploaded image to a file (you might want to store it in a more robust way)
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))
	filePath := filepath.Join("uploads/", fileName) // Adjust the path based on your requirements

	// Write the file
	newFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer func(newFile *os.File) {
		err = newFile.Close()
		if err != nil {
			return
		}
	}(newFile)

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", fmt.Errorf("error copying file: %v", err)
	}

	// Optionally, you can use the filePath or perform additional processing with the image
	fmt.Printf("Image saved: %s\n", filePath)

	return filePath, nil
}
