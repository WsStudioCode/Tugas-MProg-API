package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}
	fmt.Println(err)

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = sql.Open("mysql", conn)
	if err != nil {
		log.Fatalf("Error Connecting to Database: %v", err)
	}

	log.Println("Connected to the database successfully!")
}
