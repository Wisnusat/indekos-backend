package models

type Pembayaran struct {
	IDPembayaran     int     `json:"id_pembayaran"`
	IDReservasi      int     `json:"id_reservasi"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	JumlahBayar      float64 `json:"jumlah_bayar"`
	TanggalBayar     string  `json:"tanggal_bayar"`
	StatusPembayaran string  `json:"status_pembayaran"`
}
