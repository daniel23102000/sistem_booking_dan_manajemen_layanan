package main

import (
	"database/sql"
	"log"
	"net/http"
	"booking_system_app/database"   // Pastikan path ini sesuai dengan struktur project Anda
	"booking_system_app/middleware" // Import middleware
	_ "github.com/go-sql-driver/mysql" // Driver MySQL
)

func main() {
	// Konfigurasi DSN (Data Source Name) untuk MySQL
	dsn := "root:@tcp(127.0.0.1:3306)/booking_system" // Ganti dengan credential database Anda
	// Koneksi ke database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Menyiapkan route untuk Register dan Login
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			database.RegisterUser(db, w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			database.LoginUser(db, w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	// Menambahkan route untuk properti dan kamar dengan middleware untuk proteksi
	http.HandleFunc("/add_property", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			database.AddProperty(db, w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/add_room", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			database.AddRoom(db, w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/update_room_status", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			database.UpdateRoomStatus(db, w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}))

	// Menjalankan server HTTP di port 8080
	log.Println("Server running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}