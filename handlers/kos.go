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

func GetKosList(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id_kos, id_pemilik, nama_kos, alamat_kos, harga_sewa, deskripsi_kos, fasilitas, status_kos FROM kos")
	if err != nil {
		http.Error(w, "Failed to fetch kos data", http.StatusInternalServerError)
		log.Printf("Error fetching kos data: %v\n", err)
		return
	}
	defer rows.Close()

	var kosList []models.Kos
	for rows.Next() {
		var kos models.Kos
		err := rows.Scan(&kos.IDKos, &kos.IDPemilik, &kos.NamaKos, &kos.AlamatKos, &kos.HargaSewa, &kos.Deskripsi, &kos.Fasilitas, &kos.StatusKos)
		if err != nil {
			http.Error(w, "Failed to scan kos data", http.StatusInternalServerError)
			log.Printf("Error scanning kos data: %v\n", err)
			return
		}
		kosList = append(kosList, kos)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error occurred while processing kos data", http.StatusInternalServerError)
		log.Printf("Error processing kos data: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kosList)
}

func GetKosID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idKos := vars["id"]

	id, err := strconv.Atoi(idKos)
	if err != nil {
		http.Error(w, "Invalid kos ID", http.StatusBadRequest)
		log.Printf("Invalid kos ID: %v\n", err)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var kos models.Kos
	err = db.QueryRow(`
		SELECT id_kos, id_pemilik, nama_kos, alamat_kos, harga_sewa, deskripsi_kos, fasilitas, status_kos 
		FROM kos 
		WHERE id_kos = ?`, id).Scan(
		&kos.IDKos,
		&kos.IDPemilik,
		&kos.NamaKos,
		&kos.AlamatKos,
		&kos.HargaSewa,
		&kos.Deskripsi,
		&kos.Fasilitas,
		&kos.StatusKos,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Kos not found", http.StatusNotFound)
			log.Printf("Kos not found: ID %d\n", id)
		} else {
			http.Error(w, "Failed to fetch kos data", http.StatusInternalServerError)
			log.Printf("Error fetching kos data: %v\n", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kos)
}

func TambahKos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var kos models.Kos
	err := json.NewDecoder(r.Body).Decode(&kos)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if kos.IDPemilik == 0 || kos.NamaKos == "" || kos.AlamatKos == "" || kos.HargaSewa == 0 {
		http.Error(w, "IDPemilik, NamaKos, AlamatKos, dan HargaSewa diperlukan", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `INSERT INTO kos (id_pemilik, nama_kos, alamat_kos, harga_sewa, deskripsi, fasilitas, status_kos)
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, kos.IDPemilik, kos.NamaKos, kos.AlamatKos, kos.HargaSewa, kos.Deskripsi, kos.Fasilitas, kos.StatusKos)
	if err != nil {
		http.Error(w, "Gagal menambah data kos", http.StatusInternalServerError)
		log.Printf("Kesalahan menambah data kos: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data kos berhasil ditambahkan",
	})
}

func HapusKos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := r.URL.Query()
	idKos, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, "Invalid kos ID", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `DELETE FROM kos WHERE id_kos = ?`

	_, err = db.Exec(query, idKos)
	if err != nil {
		http.Error(w, "Gagal menghapus data kos", http.StatusInternalServerError)
		log.Printf("Kesalahan menghapus data kos: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data kos berhasil dihapus",
	})
}

func UpdateKos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var kos models.Kos
	err := json.NewDecoder(r.Body).Decode(&kos)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if kos.IDKos == 0 || kos.IDPemilik == 0 || kos.NamaKos == "" || kos.AlamatKos == "" || kos.HargaSewa == 0 {
		http.Error(w, "IDKos, IDPemilik, NamaKos, AlamatKos, dan HargaSewa diperlukan", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `UPDATE kos SET id_pemilik = ?, nama_kos = ?, alamat_kos = ?, harga_sewa = ?, deskripsi = ?, fasilitas = ?, status_kos = ? WHERE id_kos = ?`

	_, err = db.Exec(query, kos.IDPemilik, kos.NamaKos, kos.AlamatKos, kos.HargaSewa, kos.Deskripsi, kos.Fasilitas, kos.StatusKos, kos.IDKos)
	if err != nil {
		http.Error(w, "Gagal memperbarui data kos", http.StatusInternalServerError)
		log.Printf("Kesalahan memperbarui data kos: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data kos berhasil diperbarui",
	})
}
