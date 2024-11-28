package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetProfilePenyewa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var penyewa models.Penyewa
	err = db.QueryRow("SELECT id_penyewa, nama, email, no_telepon, alamat FROM penyewa WHERE id_penyewa = ?", id).Scan(
		&penyewa.IDPenyewa, &penyewa.Nama, &penyewa.Email, &penyewa.NoTelepon, &penyewa.Alamat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
			log.Printf("Error fetching user profile: %v\n", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(penyewa)
}
