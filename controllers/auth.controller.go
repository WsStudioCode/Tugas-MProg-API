package controllers

import (
	"api/config"
	"api/models"
	"api/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	var storedPassword string
	err := config.DB.QueryRow("SELECT password FROM users WHERE email = ?", user.Email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"message":"Invalid email or password"}`, http.StatusUnauthorized)
		} else {
			http.Error(w, `{"message":"Database error"}`, http.StatusInternalServerError)
			log.Println("Database error:", err)
		}
		return
	}

	if !utils.CheckPasswordHash(user.Password, storedPassword) {
		http.Error(w, `{"message":"Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	var username, email string
	var alamat sql.NullString
	err = config.DB.QueryRow("SELECT username, email, alamat FROM users WHERE email = ?", user.Email).Scan(&username, &email, &alamat)
	if err != nil {
		http.Error(w, `{"message":"Error retrieving user details"}`, http.StatusInternalServerError)
		log.Println("Error retrieving user details:", err)
		return
	}

	alamatValue := ""
	if alamat.Valid {
		alamatValue = alamat.String
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"username": username,
			"email":    email,
			"alamat":   alamatValue,
		},
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		log.Println("Error hashing password:", err)
		return
	}

	user.ID = uuid.New().String()

	_, err = config.DB.Exec("INSERT INTO users (id, username, password, email) VALUES (?,?,?,?)",
		user.ID, user.Username, hashedPassword, user.Email)

	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		log.Println("Error registering user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"success": true,
		"message": "Register successful",
	})
}
