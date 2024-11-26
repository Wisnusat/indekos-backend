package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	// Config
	user := "avnadmin"
	password := "AVNS_mn-wudrQdMHNiDuiPdP"
	host := "rental-kos-madebywisnu-19a3.b.aivencloud.com"
	port := 28725
	database := "defaultdb"

	// konfigurasi TLS
	err := registerCustomTLSConfig()
	if err != nil {
		log.Fatal("Gagal mendaftarkan konfigurasi TLS: ", err)
	}

	// Format DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?tls=custom",
		user, password, host, port, database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Koneksi ke database gagal: ", err)
	}

	// Ping
	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak dapat diakses: ", err)
	}

	log.Println("Koneksi ke database berhasil")
	return db
}

func registerCustomTLSConfig() error {
	customTLSConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	return mysql.RegisterTLSConfig("custom", customTLSConfig)
}
