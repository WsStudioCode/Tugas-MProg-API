package main

import (
	"api/config"
	"api/routes"
	"fmt"
	"log"
	"net/http"
	"os"

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
	routes.RegisterRoutes()

	fmt.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
