package models

type VerifikasiPenyewa struct {
	IDVerifikasi      int    `json:"id_verifikasi"`
	IDReservasi       int    `json:"id_reservasi"`
	StatusVerifikasi  string `json:"status_verifikasi"`
	TanggalVerifikasi string `json:"tanggal_verifikasi"`
}
