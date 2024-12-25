package routes

import (
	"api/controllers"
	"encoding/json"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "API Siap Digunakan!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CheckHealtRoutes() {
	http.HandleFunc("/health", HealthCheck)
}

func RegisterRoutes() {
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/transactions", controllers.GetAllDataTransactions)
	http.HandleFunc("/transactions/get-data", controllers.GetDataTransactionById)
	http.HandleFunc("/transactions/create", controllers.CreateTransaction)
	http.HandleFunc("/transactions/update", controllers.UpdateTransaction)
	http.HandleFunc("/transactions/delete", controllers.DeleteTransaction)
}
