package main

import (
	"backend/database"
	"backend/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// CORS
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
	db := database.ConnectDB()
	defer db.Close()

	// Router
	r := mux.NewRouter()

	r.HandleFunc("/login-pemilik-kos", handlers.LoginPemilik).Methods("POST")
	r.HandleFunc("/login-penyewa", handlers.LoginPenyewa).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")

	corsRouter := corsMiddleware(r)

	http.Handle("/", corsRouter)
	http.ListenAndServe(":8080", nil)
}
