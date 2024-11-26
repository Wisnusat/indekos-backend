package models

type PemilikKos struct {
	IDPemilik    int    `json:"id_pemilik"`
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	NomorTelepon string `json:"nomor_telepon"`
	Alamat       string `json:"alamat"`
}
