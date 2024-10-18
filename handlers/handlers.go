package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"gophercon/db" 
)

func CreateDataHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data db.Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if err := database.Create(&data).Error; err != nil {
			http.Error(w, "Failed to create data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	}
}

func ReadDataHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		column2Str := r.URL.Query().Get("column2")
		column2, err := strconv.Atoi(column2Str)
		if err != nil {
			http.Error(w, "Invalid column2 parameter", http.StatusBadRequest)
			return
		}

		var data db.Data
		if err := database.First(&data, "column2 = ?", column2).Error; err != nil {
			http.Error(w, "Data not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(data)
	}
}
