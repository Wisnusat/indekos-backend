package controllers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Reservasi oleh penyewa

// Reservasi oleh penyewa
func TambahReservasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
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

// Pemilik kos memberi persetujuan
func ApproveReservasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := r.URL.Query()
	idReservasi, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	query :=
		`UPDATE reservasi
		 SET status_reservasi = 'Diterima'
		 WHERE id_reservasi = ?;`

	_, err = db.Exec(query, idReservasi)
	if err != nil {
		http.Error(w, "Gagal menyetujui reservasi", http.StatusInternalServerError)
		log.Println("Gagal mengupdate reservasi", err)
		return

	}

	// Kirim notif ke penyewa
	w.Header().Set("Coontent-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reservasi disetujui sukses",
	})
}

// simulasi kirim notifikasi ke penyewa
func sendNotificationToPenyewa(idReservasi int, db *sql.DB) {
	query := `
		SELECT penyewa.nama, penyewa.email, penyewa.no_telepon
		FROM reservasi
		INNER JOIN penyewa ON reservasi.id_penyewa = penyewa.id_penyewa
		WHERE reservasi.id_reservasi = ?;
	`
	var penyewa models.Penyewa
	err := db.QueryRow(query, idReservasi).Scan(&penyewa.Nama, &penyewa.Email, &penyewa.NoTelepon)
	if err != nil {
		log.Println("Terjadi kesalahan saat mengambil detail penyewa ", err)
		return
	}

	log.Printf("Pesan ke penyewa [%s]: Reservasi Anda telah distujui", penyewa.Nama)
}
