package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"pencatatan-keuangan-api/config"
	"pencatatan-keuangan-api/models"
	"pencatatan-keuangan-api/utils"
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

	var username, email, alamat string
	err = config.DB.QueryRow("SELECT username, email, alamat FROM users WHERE email = ?", user.Email).Scan(&username, &email, &alamat)
	if err != nil {
		http.Error(w, `{"message":"Error retrieving user details"}`, http.StatusInternalServerError)
		log.Println("Error retrieving user details:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Login successful",
		"data": map[string]interface{}{
			"username": username,
			"email":    email,
			"alamat":   alamat,
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

	_, err = config.DB.Exec("INSERT INTO users (username, password, email, alamat) VALUES (?,?,?,?)",
		user.Username, hashedPassword, user.Email, user.Alamat)

	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		log.Println("Error registering user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Register successful",
	})
}
