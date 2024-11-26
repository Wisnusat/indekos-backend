package models

type ModerasiKos struct {
	IDModerasi      int    `json:"id_moderasi"`
	IDKos           int    `json:"id_kos"`
	StatusModerasi  string `json:"status_moderasi"`
	PesanAdmin      string `json:"pesan_admin"`
	TanggalModerasi string `json:"tanggal_moderasi"`
}
