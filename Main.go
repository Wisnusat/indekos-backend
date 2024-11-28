package main

import (
	"backend/controllers"
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

	// user
	r.HandleFunc("/login-pemilik-kos", handlers.LoginPemilik).Methods("POST")
	r.HandleFunc("/login-penyewa", handlers.LoginPenyewa).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("POST")
	r.HandleFunc("/profile-penyewa/{id}", handlers.GetProfilePenyewa).Methods("GET")
	r.HandleFunc("/profile-pemilik-kos/{id}", handlers.GetProfilePemilikKos).Methods("GET")

	// kos
	r.HandleFunc("/kos", handlers.GetKosList).Methods("GET")
	r.HandleFunc("/kos", handlers.TambahKos).Methods("POST")
	r.HandleFunc("/kos/{id}", handlers.UpdateKos).Methods("PUT")
	r.HandleFunc("/kos/{id}", handlers.HapusKos).Methods("DELETE")

	r.HandleFunc("/reservasi", controllers.TambahReservasi).Methods("POST")
	r.HandleFunc("/reservasi/approve", controllers.ApproveReservasi).Methods("PUT")

	// reservasi
	r.HandleFunc("/reservasi", handlers.TambahReservasi).Methods("POST")
	r.HandleFunc("/reservasi", handlers.GetReservasiList).Methods("GET")
	r.HandleFunc("/reservasi/{id}", handlers.UpdateReservasi).Methods("PUT")
	r.HandleFunc("/reservasi/{id}", handlers.HapusReservasi).Methods("DELETE")

	//  pembayaran
	r.HandleFunc("/pembayaran", handlers.BuatPembayaran).Methods("POST")
	r.HandleFunc("/pembayaran", handlers.GetPembayaranList).Methods("GET")
	r.HandleFunc("/pembayaran/{id}", handlers.GetPembayaranByID).Methods("GET")
	r.HandleFunc("/pembayaran/{id}", handlers.UpdatePembayaran).Methods("PUT")
	r.HandleFunc("/pembayaran/{id}", handlers.HapusPembayaran).Methods("DELETE")

	corsRouter := corsMiddleware(r)

	http.Handle("/", corsRouter)
	http.ListenAndServe(":8080", nil)
}
