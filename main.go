package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	dbp "gophercon/db"
	"gophercon/handlers" // Adjust the import path according to your project structure

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")

	dbLogger := lg.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		lg.Config{
			LogLevel:                  lg.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
			SlowThreshold:             200 * time.Millisecond,
		},
	)
	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{Logger: dbLogger})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("connected")

	// Migrate the schema for the Data model
	err = database.AutoMigrate(&dbp.Data{})
	if err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}

	database = database.Debug()

	http.HandleFunc("/create", handlers.CreateDataHandler(database))
	http.HandleFunc("/read", handlers.ReadDataHandler(database))

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
