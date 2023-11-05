package main

import (
	"Zocker-Assignment/controllers/productController"
	"Zocker-Assignment/controllers/userController"
	"Zocker-Assignment/db"
	"Zocker-Assignment/models"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"github.com/robfig/cron/v3"
)

// ProductDatabase interface for handling product database operations
type ProductDatabase interface {
	FetchProductsByImageOrderAndIsNotCompressed(limit int) ([]models.Product, error)
	UpdateProduct(updatedProduct models.Product) error
	// will Add other methods as needed
}

func main() {

	r := mux.NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello from Server !!!")
		if err != nil {
			return
		}
	}

	// Initialize MongoDB connection
	err := db.Init()
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	// Initialize your MongoDB client
	mongoClient := db.GetClient()
	userDB := db.NewMongoDBUser(mongoClient)
	productDB := db.NewMongoDB(mongoClient)

	// Initialize MongoDB collections
	err = userDB.InitUserCollection()
	if err != nil {
		fmt.Println("Failed to initialize user collections:", err)
		return
	}
	err = productDB.InitProductCollection()
	if err != nil {
		fmt.Println("Failed to initialize product collections:", err)
		return
	}

	// Image Compression logic
	go compressImagesAndUpdateProducts(productDB)

	r.HandleFunc("/", handler)
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		userController.CreateUser(w, r, userDB)
	}).Methods("POST")

	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		// Call CreateProduct with the initialized databases
		productController.CreateProduct(w, r, productDB, userDB)
	}).Methods("POST")

	port := 8080
	fmt.Printf("Server is listening on :%d...\n", port)
	http.Handle("/", r)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func compressImagesAndUpdateProducts(productDB ProductDatabase) {
	// Use a WaitGroup to wait for all goroutines to finish.
	var mainWg sync.WaitGroup

	// Run the cron job in a goroutine.
	mainWg.Add(1)

	go func() {
		defer mainWg.Done()

		// Create a cron-like scheduler.
		c := cron.New()

		// Schedule the cron job to run every 5 minutes.
		_, err := c.AddFunc("*/1 * * * *", func() {
			// Fetch products with uncompressed images (limit 10).
			uncompressedProducts, err := productDB.FetchProductsByImageOrderAndIsNotCompressed(10)
			if err != nil {
				log.Fatalf("Error fetching products: %v", err)
			}
			// Use a WaitGroup to wait for all compression goroutines to finish.
			var wg sync.WaitGroup

			// Iterate over uncompressed products and start a goroutine for each product.
			for _, product := range uncompressedProducts {
				wg.Add(1)
				go CompressImagesForProductAndUpdateProduct(product, &wg, productDB)
			}
			// Wait for all product compression goroutines to finish.
			wg.Wait()
		})

		if err != nil {
			log.Fatal("Error scheduling cron job:", err)
		}

		// Start the cron scheduler.
		c.Start()

		// Keep the program running to allow the cron job to execute.
		select {}
	}()
}

// CompressImagesForProductAndUpdateProduct compresses images for a given product
// and updates the product with the compressed image paths.
func CompressImagesForProductAndUpdateProduct(product models.Product, wg *sync.WaitGroup, productDB ProductDatabase) {
	defer wg.Done()

	// Use a WaitGroup to wait for all compression goroutines to finish.
	var imageWg sync.WaitGroup

	// Create a channel to receive compressed image paths and errors.
	compressedImageChan := make(chan string, len(product.ProductImages))
	errorChan := make(chan error, len(product.ProductImages))
	var firstError error

	// Iterate over product images and start a goroutine for each image.
	for _, productImage := range product.ProductImages {
		imageWg.Add(1)
		go func(productImage string) {
			defer imageWg.Done()
			// Compress the image.
			compressedImage, err := CompressImage(productImage)
			if err != nil {
				// Handle compression error.
				fmt.Printf("Error compressing image %s: %v\n", productImage, err)
				errorChan <- err
				return
			}

			// Send the compressed image path to the channel.
			compressedImageChan <- compressedImage
		}(productImage)
	}

	// Close the channels when all image compression goroutines are done.
	go func() {
		imageWg.Wait()
		close(compressedImageChan)
		close(errorChan)
	}()

	// Range over the error channel to handle errors.
	for err := range errorChan {
		fmt.Printf("Compression error: %v\n", err)
		firstError = err
	}

	if firstError != nil {
		return
	}

	// Collect compressed image paths and handle errors.
	var compressedImages []string
	for compressedImage := range compressedImageChan {
		compressedImages = append(compressedImages, compressedImage)
	}

	// Update the product with the compressed image paths.
	// In a real scenario, you would update the product in the database.
	// Here, we're just updating the local product for demonstration purposes.
	product.IsCompressed = true
	product.CompressedProductImages = compressedImages

	err := productDB.UpdateProduct(product)
	if err != nil {
		fmt.Printf("error updating the product, %s", err)
	}
	// Print the result for demonstration purposes.
	fmt.Printf("Product %d compressed successfully. Compressed images: %v\n", product.ID, compressedImages)
}

func CompressImage(inputImagePath string) (string, error) {
	maxWidth := uint(800)
	maxHeight := uint(600)

	outputImagePath := fmt.Sprintf("compressedImageUploads/%s", filepath.Base(inputImagePath))

	// Open the image file
	file, err := os.Open(inputImagePath)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", inputImagePath, err)
		return "", fmt.Errorf("error opening input file %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding image %s: %v\n", inputImagePath, err)
		return "", fmt.Errorf("error decoding input image %w", err)
	}

	// Resize the image
	resized := resize.Resize(maxWidth, maxHeight, img, resize.Lanczos3)

	// Create the output file
	outputFile, err := os.Create(outputImagePath)
	if err != nil {
		log.Printf("Error creating output file %s: %v\n", outputImagePath, err)
		return "", fmt.Errorf("error creating output file %w", err)
	}
	defer func(outFile *os.File) {
		err = outFile.Close()
		if err != nil {
			return
		}
	}(outputFile)

	// Encode the resized image to the output file
	err = jpeg.Encode(outputFile, resized, nil)
	if err != nil {
		log.Printf("Error encoding image to %s: %v\n", outputImagePath, err)
		return "", fmt.Errorf("error encoding image to output image path %w", err)
	}

	fmt.Printf("Image %s compressed successfully!\n", inputImagePath)
	return outputImagePath, nil
}
