package controllers

import (
	"api/config"
	"api/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func GetAllDataTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []models.Transaction
	result, err := config.DB.Query("SELECT id, amount, description, date, userId, category FROM transactions")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	for result.Next() {
		var transaction models.Transaction
		err := result.Scan(&transaction.ID, &transaction.Amount, &transaction.Description, &transaction.Date, &transaction.UserID, &transaction.Category)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := models.APIResponse{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "No transactions found",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.APIResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Data transactions retrieved successfully",
		Data:    transactions,
	}

	log.Printf("Response: %+v", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDataTransactionById(w http.ResponseWriter, r *http.Request) {
	var transactions []models.Transaction

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, `{"message":"Missing transaction ID"}`, http.StatusBadRequest)
		return
	}

	query := "SELECT id, amount, description, date, userId, category FROM transactions WHERE id =?"
	rows, err := config.DB.Query(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.Description,
			&transaction.Date,
			&transaction.UserID,
			&transaction.Category,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := models.APIResponse{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "No transactions found with given ID",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.APIResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Data transactions retrieved successfully by id",
		Data:    transactions,
	}

	log.Printf("Response: %+v", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	transaction.ID = uuid.New().String()

	query := "INSERT INTO transactions (id, amount, description, date, userId, category) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = config.DB.Exec(query, transaction.ID, transaction.Amount, transaction.Description, transaction.Date, transaction.UserID, transaction.Category)
	if err != nil {
		http.Error(w, `{"message":"Failed to create transaction: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := models.APIResponse{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Successfull created data transaction",
		Data:    transaction,
	}

	log.Printf("Response: %+v", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"message":"Missing transaction ID"}`, http.StatusBadRequest)
		return
	}

	var transaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	query := "UPDATE transactions SET amount = ?, description = ?, date = ?, userId = ?, category = ? WHERE id = ?"
	_, err = config.DB.Exec(query, transaction.Amount, transaction.Description, transaction.Date, transaction.UserID, transaction.Category, id)
	if err != nil {
		http.Error(w, `{"message":"Failed to update transaction: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := models.APIResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Transaction updated successfully",
		Data:    transaction,
	}

	log.Printf("Response: %+v", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"message":"Missing transaction ID"}`, http.StatusBadRequest)
		return
	}
	query := "DELETE FROM transactions WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, `{"message":"Failed to delete transaction: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := models.APIResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Transaction deleted successfully",
	}

	log.Printf("Response: %+v", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
