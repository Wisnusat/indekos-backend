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

func GetProfilePemilikKos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var pemilik models.PemilikKos
	err = db.QueryRow("SELECT id_pemilik, nama, email, nomor_telepon, alamat FROM pemilik_kos WHERE id_pemilik = ?", id).Scan(
		&pemilik.IDPemilik, &pemilik.Nama, &pemilik.Email, &pemilik.NomorTelepon, &pemilik.Alamat)
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
	json.NewEncoder(w).Encode(pemilik)
}
