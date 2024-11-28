package handlers

import (
	"backend/database"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func TambahReservasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reservasiReq struct {
		IDPenyewa int `json:"id_penyewa"`
		IDKos     int `json:"id_kos"`
	}

	err := json.NewDecoder(r.Body).Decode(&reservasiReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reservasiReq.IDPenyewa == 0 || reservasiReq.IDKos == 0 {
		http.Error(w, "IDPenyewa and IDKos are required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	// Validasi IDKos dan IDPenyewa
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM kos WHERE id_kos = ?", reservasiReq.IDKos).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "IDKos tidak valid", http.StatusBadRequest)
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM penyewa WHERE id_penyewa = ?", reservasiReq.IDPenyewa).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "IDPenyewa tidak valid", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO reservasi (id_penyewa, id_kos, tanggal_reservasi, status_reservasi)
              VALUES (?, ?, ?, ?);`

	tanggalReservasi := time.Now().Format("2006-01-02")
	statusReservasi := "Menunggu Verifikasi"

	_, err = db.Exec(query, reservasiReq.IDPenyewa, reservasiReq.IDKos, tanggalReservasi, statusReservasi)
	if err != nil {
		http.Error(w, "Gagal membuat reservasi", http.StatusInternalServerError)
		log.Printf("Kesalahan memasukkan reservasi: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reservasi berhasil dibuat, menunggu persetujuan",
	})
}

func GetReservasiList(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id_reservasi, id_penyewa, id_kos, tanggal_reservasi, status_reservasi FROM reservasi")
	if err != nil {
		http.Error(w, "Failed to fetch reservasi data", http.StatusInternalServerError)
		log.Printf("Error fetching reservasi data: %v\n", err)
		return
	}
	defer rows.Close()

	var reservasiList []models.Reservasi
	for rows.Next() {
		var reservasi models.Reservasi
		err := rows.Scan(&reservasi.IDReservasi, &reservasi.IDPenyewa, &reservasi.IDKos, &reservasi.TanggalReservasi, &reservasi.StatusReservasi)
		if err != nil {
			http.Error(w, "Failed to scan reservasi data", http.StatusInternalServerError)
			log.Printf("Error scanning reservasi data: %v\n", err)
			return
		}
		reservasiList = append(reservasiList, reservasi)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error occurred while processing reservasi data", http.StatusInternalServerError)
		log.Printf("Error processing reservasi data: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservasiList)
}

func UpdateReservasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reservasi models.Reservasi
	err := json.NewDecoder(r.Body).Decode(&reservasi)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reservasi.IDReservasi == 0 || reservasi.IDPenyewa == 0 || reservasi.IDKos == 0 {
		http.Error(w, "IDReservasi, IDPenyewa, dan IDKos diperlukan", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `UPDATE reservasi SET id_penyewa = ?, id_kos = ?, tanggal_reservasi = ?, status_reservasi = ? WHERE id_reservasi = ?`

	_, err = db.Exec(query, reservasi.IDPenyewa, reservasi.IDKos, reservasi.TanggalReservasi, reservasi.StatusReservasi, reservasi.IDReservasi)
	if err != nil {
		http.Error(w, "Gagal memperbarui data reservasi", http.StatusInternalServerError)
		log.Printf("Kesalahan memperbarui data reservasi: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data reservasi berhasil diperbarui",
	})
}

func HapusReservasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := r.URL.Query()
	idReservasi, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, "Invalid reservasi ID", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `DELETE FROM reservasi WHERE id_reservasi = ?`

	_, err = db.Exec(query, idReservasi)
	if err != nil {
		http.Error(w, "Gagal menghapus data reservasi", http.StatusInternalServerError)
		log.Printf("Kesalahan menghapus data reservasi: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data reservasi berhasil dihapus",
	})
}
