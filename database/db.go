package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDB() *sql.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || database == "" {
		log.Fatal("Database configuration is missing in environment variables")
	}

	// TLS configuration
	err = registerCustomTLSConfig()
	if err != nil {
		log.Fatal("Failed to register custom TLS configuration: ", err)
	}

	// Format DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=custom",
		user, password, host, port, database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Ping  database
	err = db.Ping()
	if err != nil {
		log.Fatal("Database is not accessible: ", err)
	}

	log.Println("Connected to the database successfully")
	return db
}

func registerCustomTLSConfig() error {
	customTLSConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	return mysql.RegisterTLSConfig("custom", customTLSConfig)
}
