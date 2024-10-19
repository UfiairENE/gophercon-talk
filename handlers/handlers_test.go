package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	dbp "gophercon/db"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, func()) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(dbp.Data{})
	if err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}

	return db, func() {
		// db.Exec("DELETE FROM data") // Cleanup after tests
	}
}

func TestCreateDataHandler(t *testing.T) {
	db, teardown := setupTestDB()
	defer teardown()

	handler := CreateDataHandler(db)

	t.Run("Valid request", func(t *testing.T) {
		data := dbp.Data{Column1: "Test", Column2: 1}
		body, _ := json.Marshal(data)

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		var responseData dbp.Data
		json.NewDecoder(rec.Body).Decode(&responseData)
		assert.Equal(t, data.Column1, responseData.Column1)
		assert.Equal(t, data.Column2, responseData.Column2)
	})

	t.Run("Invalid method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/create", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Invalid payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte("invalid")))
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestReadDataHandler(t *testing.T) {
	db, teardown := setupTestDB()
	defer teardown()

	col2 := 200000
	db.Exec("delete from data where column2=?", col2)
	data := dbp.Data{Column1: "Test", Column2: col2}
	db.Create(&data)

	handler := ReadDataHandler(db)

	t.Run("Valid request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/read?column2=%v", col2), nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var responseData dbp.Data
		json.NewDecoder(rec.Body).Decode(&responseData)
		assert.Equal(t, data.Column1, responseData.Column1)
		assert.Equal(t, data.Column2, responseData.Column2)
	})

	t.Run("Data not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/read?column2=99999999999999", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
