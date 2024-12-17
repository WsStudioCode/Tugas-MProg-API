package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pencatatan-keuangan-api/config"
	"pencatatan-keuangan-api/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading.env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	config.ConnectDB()
	routes.RegisterAuthRoutes()

	fmt.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
