package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type PaymentRequest struct {
	IDReservasi      int     `json:"id_reservasi"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	JumlahBayar      float64 `json:"jumlah_bayar"`
}

func BuatPembayaran(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var paymentReq PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if paymentReq.IDReservasi == 0 || paymentReq.MetodePembayaran == "" || paymentReq.JumlahBayar == 0 {
		http.Error(w, "IDReservasi, MetodePembayaran, dan JumlahBayar diperlukan", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	// Validasi IDReservasi
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM reservasi WHERE id_reservasi = ?", paymentReq.IDReservasi).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "IDReservasi tidak valid", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO pembayaran (id_reservasi, metode_pembayaran, jumlah_bayar, tanggal_bayar, status_pembayaran)
              VALUES (?, ?, ?, ?, ?);`

	tanggalBayar := time.Now().Format("2006-01-02")
	statusPembayaran := "Sudah Dibayar"

	_, err = db.Exec(query, paymentReq.IDReservasi, paymentReq.MetodePembayaran, paymentReq.JumlahBayar, tanggalBayar, statusPembayaran)
	if err != nil {
		http.Error(w, "Gagal memproses pembayaran", http.StatusInternalServerError)
		log.Printf("Kesalahan memproses pembayaran: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Pembayaran berhasil diproses",
	})
}

func GetPembayaranList(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id_pembayaran, id_reservasi, metode_pembayaran, jumlah_bayar, tanggal_bayar, status_pembayaran FROM pembayaran")
	if err != nil {
		http.Error(w, "Failed to fetch pembayaran data", http.StatusInternalServerError)
		log.Printf("Error fetching pembayaran data: %v\n", err)
		return
	}
	defer rows.Close()

	var pembayaranList []models.Pembayaran
	for rows.Next() {
		var pembayaran models.Pembayaran
		err := rows.Scan(&pembayaran.IDPembayaran, &pembayaran.IDReservasi, &pembayaran.MetodePembayaran, &pembayaran.JumlahBayar, &pembayaran.TanggalBayar, &pembayaran.StatusPembayaran)
		if err != nil {
			http.Error(w, "Failed to scan pembayaran data", http.StatusInternalServerError)
			log.Printf("Error scanning pembayaran data: %v\n", err)
			return
		}
		pembayaranList = append(pembayaranList, pembayaran)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error occurred while processing pembayaran data", http.StatusInternalServerError)
		log.Printf("Error processing pembayaran data: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pembayaranList)
}

func GetPembayaranByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPembayaran, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pembayaran ID", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var pembayaran models.Pembayaran
	query := `SELECT id_pembayaran, id_reservasi, metode_pembayaran, jumlah_bayar, tanggal_bayar, status_pembayaran FROM pembayaran WHERE id_pembayaran = ?`
	err = db.QueryRow(query, idPembayaran).Scan(&pembayaran.IDPembayaran, &pembayaran.IDReservasi, &pembayaran.MetodePembayaran, &pembayaran.JumlahBayar, &pembayaran.TanggalBayar, &pembayaran.StatusPembayaran)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Pembayaran tidak ditemukan", http.StatusNotFound)
		} else {
			http.Error(w, "Gagal mengambil data pembayaran", http.StatusInternalServerError)
			log.Printf("Kesalahan mengambil data pembayaran: %v\n", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pembayaran)
}

func UpdatePembayaran(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var pembayaran models.Pembayaran
	err := json.NewDecoder(r.Body).Decode(&pembayaran)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if pembayaran.IDPembayaran == 0 || pembayaran.IDReservasi == 0 || pembayaran.MetodePembayaran == "" || pembayaran.JumlahBayar == 0 {
		http.Error(w, "IDPembayaran, IDReservasi, MetodePembayaran, dan JumlahBayar diperlukan", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `UPDATE pembayaran SET id_reservasi = ?, metode_pembayaran = ?, jumlah_bayar = ?, tanggal_bayar = ?, status_pembayaran = ? WHERE id_pembayaran = ?`

	_, err = db.Exec(query, pembayaran.IDReservasi, pembayaran.MetodePembayaran, pembayaran.JumlahBayar, pembayaran.TanggalBayar, pembayaran.StatusPembayaran, pembayaran.IDPembayaran)
	if err != nil {
		http.Error(w, "Gagal memperbarui data pembayaran", http.StatusInternalServerError)
		log.Printf("Kesalahan memperbarui data pembayaran: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data pembayaran berhasil diperbarui",
	})
}

func HapusPembayaran(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := r.URL.Query()
	idPembayaran, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, "Invalid pembayaran ID", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query := `DELETE FROM pembayaran WHERE id_pembayaran = ?`

	_, err = db.Exec(query, idPembayaran)
	if err != nil {
		http.Error(w, "Gagal menghapus data pembayaran", http.StatusInternalServerError)
		log.Printf("Kesalahan menghapus data pembayaran: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data pembayaran berhasil dihapus",
	})
}
