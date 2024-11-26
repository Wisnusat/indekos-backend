package models

type Reservasi struct {
	IDReservasi      int    `json:"id_reservasi"`
	IDPenyewa        int    `json:"id_penyewa"`
	IDKos            int    `json:"id_kos"`
	TanggalReservasi string `json:"tanggal_reservasi"`
	StatusReservasi  string `json:"status_reservasi"`
}
