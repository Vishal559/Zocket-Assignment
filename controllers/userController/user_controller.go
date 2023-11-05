package userController

import (
	"Zocker-Assignment/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserDatabase interface {
	InsertUser(user models.User) error
}

func CreateUser(w http.ResponseWriter, r *http.Request, userDB UserDatabase) {
	var newUser models.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the required fields (we can add more validation as needed)
	if newUser.UserId == "" || newUser.Name == "" || newUser.Mobile == "" {
		http.Error(w, "userid, name and mobile no. are required fields", http.StatusBadRequest)
		return
	}

	// Create the user in the database
	err = userDB.InsertUser(newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the created user (you may choose to exclude sensitive fields in the response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newUser)
	if err != nil {
		return
	}
}
