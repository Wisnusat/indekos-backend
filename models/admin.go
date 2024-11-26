package models

type Admin struct {
	IDAdmin      int    `json:"id_admin"`
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	NomorTelepon string `json:"nomor_telepon"`
}
