package handlers

import (
	"backend/database"
	"encoding/json"
	"log"
	"net/http"
)

type RegisterRequest struct {
	UserType  string `json:"user_type"`
	Nama      string `json:"nama"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	NoTelepon string `json:"no_telepon"`
	Alamat    string `json:"alamat"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var registerReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if registerReq.UserType == "" || registerReq.Nama == "" || registerReq.Email == "" || registerReq.Password == "" || registerReq.NoTelepon == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if registerReq.UserType != "penyewa" && registerReq.UserType != "pemilik_kos" {
		http.Error(w, "Invalid user type", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var query string
	if registerReq.UserType == "penyewa" {
		query = `
			INSERT INTO penyewa (nama, email, password, no_telepon, alamat)
			VALUES (?, ?, ?, ?, ?);
		`
	} else if registerReq.UserType == "pemilik_kos" {
		query = `
			INSERT INTO pemilik_kos (nama, email, password, nomor_telepon, alamat)
			VALUES (?, ?, ?, ?, ?);
		`
	}

	_, err = db.Exec(query, registerReq.Nama, registerReq.Email, registerReq.Password, registerReq.NoTelepon, registerReq.Alamat)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Println("Error inserting user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		Message: "User registered successfully",
	})
}
