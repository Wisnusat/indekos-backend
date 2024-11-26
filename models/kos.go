package models

type Kos struct {
	IDKos     int     `json:"id_kos"`
	IDPemilik int     `json:"id_pemilik"`
	NamaKos   string  `json:"nama_kos"`
	AlamatKos string  `json:"alamat_kos"`
	HargaSewa float64 `json:"harga_sewa"`
	Deskripsi string  `json:"deskripsi"`
	Fasilitas string  `json:"fasilitas"`
	StatusKos string  `json:"status_kos"`
}
