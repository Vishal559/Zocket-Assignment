package productController

import (
	"Zocker-Assignment/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockProductDB is a mock implementation of ProductDatabase
type MockProductDB struct {
	mock.Mock
}

func (m *MockProductDB) InitProductCollection() error {
	return nil
}

func (m *MockProductDB) InsertProduct(product models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// MockUserDB is a mock implementation of UserDatabase
type MockUserDB struct {
	mock.Mock
}

func (m *MockUserDB) FindUser(userID string) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}

// createMultipartRequest creates an HTTP request with a JSON payload
func createMultipartRequest(jsonData string, targetURL string) (*http.Request, error) {
	// Unmarshal the JSON data into a map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}

	// Create a buffer to write the payload
	var payload bytes.Buffer

	// Create a new multipart writer with the buffer
	writer := multipart.NewWriter(&payload)

	// Add fields from the map to the multipart writer
	for key, value := range data {
		_ = writer.WriteField(key, fmt.Sprintf("%v", value))
	}

	// Close the writer to finalize the payload
	_ = writer.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", targetURL, &payload)
	if err != nil {
		return nil, err
	}

	// Set the appropriate Content-Type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func TestCreateProduct(t *testing.T) {
	// Set up the test data
	testData := []struct {
		name         string
		payload      string
		expectedCode int
		setupMocks   func(productDB *MockProductDB, userDB *MockUserDB)
	}{
		// Happy flow
		{
			name:         "Valid request",
			payload:      `{"user_id": "123", "product_name": "Test Product", "product_description": "Test Description", "product_price": "50"}`,
			expectedCode: http.StatusCreated,
			setupMocks: func(productDB *MockProductDB, userDB *MockUserDB) {
				userDB.On("FindUser", mock.Anything).Return(&models.User{
					UserId: "123",
				}, nil)
				productDB.On("InsertProduct", mock.AnythingOfType("Product")).Return(nil)
			},
		},
		// Validation error
		{
			name:         "Missing required fields",
			payload:      `{"user_id": "123"}`,
			expectedCode: http.StatusBadRequest,
			setupMocks:   func(productDB *MockProductDB, userDB *MockUserDB) {},
		},
		// Database error
		{
			name:         "Database error",
			payload:      `{"user_id": "123", "product_name": "Test Product", "product_description": "Test Description", "product_price": "50"}`,
			expectedCode: http.StatusInternalServerError,
			setupMocks: func(productDB *MockProductDB, userDB *MockUserDB) {
				userDB.On("FindUser", mock.Anything).Return(&models.User{
					UserId: "123",
				}, nil)
				productDB.On("InsertProduct", mock.AnythingOfType("Product")).Return(errors.New("database error"))
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create mock databases
			mockProductDB := new(MockProductDB)
			mockUserDB := new(MockUserDB)

			// Set up expectations on the mock databases
			tt.setupMocks(mockProductDB, mockUserDB)

			// Create a new request with the mock form data
			request, err := createMultipartRequest(tt.payload, "/products")
			if err != nil {
				panic("failed to generate http request")
			}

			// Call the handler function
			CreateProduct(rr, request, mockProductDB, mockUserDB)

			// Check the response status code
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			// Check if the mock database methods were called as expected
			mockProductDB.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}
