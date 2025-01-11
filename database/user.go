package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

// Struct untuk request registrasi
type RegisterRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

// Struct untuk request login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct untuk properti
type Property struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Description  string `json:"description"`
	ContactNumber string `json:"contact_number"`
}

// Struct untuk kamar
type Room struct {
	PropertyID     int     `json:"property_id"`
	RoomName       string  `json:"room_name"`
	RoomType       string  `json:"room_type"`
	PricePerNight  float64 `json:"price_per_night"`
	Status         string  `json:"status"` // tersedia, dipesan, atau dalam perawatan
}

// Struct untuk response login dan register
type Response struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// Struct untuk JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var secretKey = []byte("your_secret_key") // Gantilah dengan kunci yang aman

// RegisterUser menangani registrasi user baru
func RegisterUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Decode data JSON dari body request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password menggunakan bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// SQL untuk memasukkan user baru
	query := `INSERT INTO users (name, email, password_hash, phone_number, role) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(query, req.Name, req.Email, hash, req.PhoneNumber, req.Role)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error registering user: %v", err), http.StatusInternalServerError)
		return
	}

	// Response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Message: "User registered successfully"})
}

// LoginUser menangani proses login user
func LoginUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Decode data JSON dari body request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// SQL untuk memvalidasi user
	query := `SELECT user_id, name, role, password_hash FROM users WHERE email = ?`
	row := db.QueryRow(query, req.Email)

	var userID int
	var name, role, storedHash string
	err = row.Scan(&userID, &name, &role, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	// Verifikasi password menggunakan bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Membuat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Could not create JWT token", http.StatusInternalServerError)
		return
	}

	// Kirimkan token sebagai respons
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Login successful. Welcome, " + name + " (Role: " + role + ")",
		Token:   tokenString,
	})
}

// AddProperty menangani penambahan properti baru
func AddProperty(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var property Property

	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO properties (name, address, description, contact_number) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, property.Name, property.Address, property.Description, property.ContactNumber)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding property: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Message: "Property added successfully"})
}

// AddRoom menangani penambahan kamar baru
func AddRoom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var room Room

	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO rooms (property_id, room_type, price_per_night, status) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, room.PropertyID, room.RoomType, room.PricePerNight, room.Status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding room: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Message: "Room added successfully"})
}

// UpdateRoomStatus menangani pembaruan status kamar
func UpdateRoomStatus(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	type UpdateStatusRequest struct {
		RoomID int    `json:"room_id"`
		Status string `json:"status"` // tersedia, dipesan, atau dalam perawatan
	}

	var req UpdateStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE rooms SET status = ? WHERE id = ?`
	_, err = db.Exec(query, req.Status, req.RoomID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating room status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "Room status updated successfully"})
}