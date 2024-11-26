package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type LoginRequest struct {
	Email     string `json:"email,omitempty"`
	NoTelepon string `json:"no_telepon,omitempty"`
	Password  string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

func LoginPemilik(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if loginReq.Email == "" && loginReq.NoTelepon == "" {
		http.Error(w, "Email or NoTelepon is required", http.StatusBadRequest)
		return
	}
	if loginReq.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `
		SELECT id_pemilik, nama, email, password, nomor_telepon, alamat
		FROM pemilik_kos
		WHERE (email = ? OR nomor_telepon = ?)
		AND password = ?;
	`

	var user models.PemilikKos
	err = db.QueryRow(query, loginReq.Email, loginReq.NoTelepon, loginReq.Password).Scan(
		&user.IDPemilik,
		&user.Nama,
		&user.Email,
		&user.Password,
		&user.NomorTelepon,
		&user.Alamat,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email/no telepon or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Println("Error scanning result: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Login successful",
	})
}

func LoginPenyewa(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if loginReq.Email == "" && loginReq.NoTelepon == "" {
		http.Error(w, "Email or NoTelepon is required", http.StatusBadRequest)
		return
	}
	if loginReq.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `
		SELECT id_penyewa, nama, email, password, no_telepon, alamat
		FROM penyewa
		WHERE (email = ? OR no_telepon = ?)
		AND password = ?;
	`

	var user models.Penyewa
	err = db.QueryRow(query, loginReq.Email, loginReq.NoTelepon, loginReq.Password).Scan(
		&user.IDPenyewa,
		&user.Nama,
		&user.Email,
		&user.Password,
		&user.NoTelepon,
		&user.Alamat,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email/no telepon or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Println("Error scanning result: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Login successful",
	})
}
