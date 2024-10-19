package main

import (
	"log"
	"os"
	"strconv"

	dbp "gophercon/db"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema for the Data model
	err = database.AutoMigrate(&dbp.Data{})
	if err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}

	dataLength := 10000000
	batchLength := 5000

	// dataLength := 100
	// batchLength := 100

	err = insertDataWithTransaction(database, dataLength, batchLength)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}

	log.Println("Data seeding completed successfully.")
}

func insertDataWithTransaction(database *gorm.DB, dataLength, batchLength int) error {
	tx := database.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction failed. Rolled back.")
		}
	}()

	dataBatches := prepareData(dataLength, batchLength)

	for _, batch := range dataBatches {
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Printf("Inserted %d records successfully.", dataLength)
	return nil
}

func prepareData(totalSize, batchLength int) [][]dbp.Data {
	var batches [][]dbp.Data
	var batch []dbp.Data

	for i := 0; i < totalSize; i++ {
		batch = append(batch, dbp.Data{
			Column1: "Value " + strconv.Itoa(i),
			Column2: i,
		})

		if len(batch) == batchLength {
			batches = append(batches, batch)
			batch = []dbp.Data{}
		}
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}
