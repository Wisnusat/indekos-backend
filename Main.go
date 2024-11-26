package main

import (
	"backend/database"
	"backend/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db := database.ConnectDB()
	defer db.Close()

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/login-pemilik-kos", handlers.LoginPemilik).Methods("POST")
	r.HandleFunc("/login-penyewa", handlers.LoginPenyewa).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	
	corsRouter := corsMiddleware(r)
	
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: corsRouter,
	}

	log.Printf("Server starting on port %s", port)
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
