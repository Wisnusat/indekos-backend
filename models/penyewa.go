package models

type Penyewa struct {
	IDPenyewa int    `json:"id_penyewa"`
	Nama      string `json:"nama"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	NoTelepon string `json:"no_telepon"`
	Alamat    string `json:"alamat"`
}
