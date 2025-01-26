package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"
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

type BookingDetail struct {
    RoomID         int      `json:"room_id"`
    Quantity      int     `json:"quantity"`
    PricePerNight float64 `json:"price_per_night"`
}

type BookingRequest struct {
    CustomerID       int                `json:"customer_id"`
    CheckInDate      string             `json:"check_in_date"`
    CheckOutDate     string             `json:"check_out_date"`
    BookingDetails   []BookingDetail    `json:"booking_details"`
    AdditionalServices []int            `json:"additional_services"`
    PaymentDetails   PaymentDetails     `json:"payment_details"`
}

type BookingDetails struct {
    RoomID   int `json:"room_id"`
    Quantity int `json:"quantity"`
}

type PaymentDetails struct {
    PaymentMethod string  `json:"payment_method"`
    TotalAmount   float64 `json:"total_amount"`
}

type BookingResponse struct {
    BookingIDs []int   `json:"booking_ids"`
    TotalPrice float64 `json:"total_price"`
    Message    string  `json:"message"`
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
	type Room struct {
		PropertyID    int     `json:"property_id"`
		RoomName      string  `json:"room_name"`
		RoomType      string  `json:"room_type"`
		PricePerNight float64 `json:"price_per_night"`
		Status        string  `json:"status"`
	}

	type Response struct {
		Message string `json:"message"`
	}

	var room Room

	// Decode body request
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log nilai yang diterima untuk debugging
	fmt.Println("Received Room Data:", room)

	// Validasi input (opsional)
	if room.PropertyID <= 0 || room.RoomName == "" || room.RoomType == "" || room.PricePerNight <= 0 || room.Status == "" {
		// Log jika ada input yang tidak valid
		fmt.Println("Invalid input data:", room)
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// Log sebelum query dieksekusi
	fmt.Printf("Prepared to insert room with PropertyID: %d, RoomName: %s, RoomType: %s, PricePerNight: %.2f, Status: %s\n",
		room.PropertyID, room.RoomName, room.RoomType, room.PricePerNight, room.Status)

	// Query untuk menambahkan kamar
	query := `
		INSERT INTO rooms (property_id, room_name, room_type, price_per_night, status)
		VALUES (?, ?, ?, ?, ?)
	`

	// Eksekusi query dan tangkap error jika ada
	fmt.Println("room.RoomType",room.RoomType)
	result, err := db.Exec(query, room.PropertyID, room.RoomName, room.RoomType, room.PricePerNight, room.Status)
	if err != nil {
		// Log error eksekusi query
		fmt.Println("Error executing query:", err)
		http.Error(w, fmt.Sprintf("Error adding room: %v", err), http.StatusInternalServerError)
		return
	}

	// Log hasil eksekusi query
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		// Log error untuk mendapatkan ID terakhir
		fmt.Println("Error getting last insert ID:", err)
		http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
		return
	}

	// Log jika data berhasil ditambahkan
	fmt.Printf("Room added successfully with ID: %d\n", lastInsertID)

	// Respons sukses
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

	query := `UPDATE rooms SET status = ? WHERE room_id = ?`
	_, err = db.Exec(query, req.Status, req.RoomID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating room status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "Room status updated successfully"})
}

// SearchRooms menangani pencarian kamar
func SearchRooms(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    // Struktur untuk menerima kriteria pencarian
    type SearchCriteria struct {
        PropertyName string  `json:"property_name"`
        RoomType     string  `json:"room_type"`
        MinPrice     float64 `json:"min_price"`
        MaxPrice     float64 `json:"max_price"`
    }

    var criteria SearchCriteria
    err := json.NewDecoder(r.Body).Decode(&criteria)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validasi harga (min_price tidak boleh lebih besar dari max_price)
    if criteria.MinPrice > criteria.MaxPrice {
        http.Error(w, "min_price cannot be greater than max_price", http.StatusBadRequest)
        return
    }

    // Gunakan wildcard jika property_name atau room_type kosong
    propertyName := "%"
    if criteria.PropertyName != "" {
        propertyName = "%" + criteria.PropertyName + "%"
    }

    roomType := "%"
    if criteria.RoomType != "" {
        roomType = "%" + criteria.RoomType + "%"
    }

    // Log nilai kriteria pencarian untuk debugging
    fmt.Printf("Searching with property_name: %s, room_type: %s, min_price: %.2f, max_price: %.2f\n",
        criteria.PropertyName, criteria.RoomType, criteria.MinPrice, criteria.MaxPrice)

    // Query SQL untuk pencarian kamar
    query := `
        SELECT r.room_id, r.room_name, r.room_type, r.price_per_night, r.status, p.name AS property_name
        FROM rooms r
        JOIN properties p ON r.property_id = p.property_id
        WHERE p.name LIKE ? AND r.room_type LIKE ? AND r.price_per_night BETWEEN ? AND ?
    `

    // Log query dan parameter untuk debugging
    fmt.Printf("Query: %s\nParams: propertyName=%s, roomType=%s, minPrice=%.2f, maxPrice=%.2f\n",
        query, propertyName, roomType, criteria.MinPrice, criteria.MaxPrice)

    // Eksekusi query dengan parameter pencarian
    rows, err := db.Query(query, propertyName, roomType, criteria.MinPrice, criteria.MaxPrice)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error Searching Rooms: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Struktur untuk hasil pencarian kamar
    type RoomSearchResult struct {
        ID            int     `json:"id"`
        RoomName      string  `json:"room_name"`
        RoomType      string  `json:"room_type"`
        PricePerNight float64 `json:"price_per_night"`
        Status        string  `json:"status"`
        PropertyName  string  `json:"property_name"`
    }

    var results []RoomSearchResult

    // Memindai hasil query dan mengisi hasil pencarian
    for rows.Next() {
        var result RoomSearchResult
        err := rows.Scan(&result.ID, &result.RoomName, &result.RoomType, &result.PricePerNight, &result.Status, &result.PropertyName)
        if err != nil {
            http.Error(w, "Error reading search results", http.StatusInternalServerError)
            return
        }
        results = append(results, result)
    }

    // Jika tidak ada hasil ditemukan
    if len(results) == 0 {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"message": "No rooms found"})
        return
    }

    // Respons sukses dengan hasil pencarian
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func BookRoom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var req BookingRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validasi input
    log.Println("req", req)
    if req.CustomerID <= 0 || len(req.BookingDetails) == 0 || req.CheckInDate == "" || req.CheckOutDate == "" || req.PaymentDetails.TotalAmount <= 0 {
        http.Error(w, fmt.Sprintf("Invalid Input Data: %v", err), http.StatusBadRequest)
        return
    }

    // Parse tanggal check-in dan check-out
    checkIn, err := time.Parse("2006-01-02", req.CheckInDate)
    checkOut, err := time.Parse("2006-01-02", req.CheckOutDate)
    if err != nil || checkOut.Before(checkIn) {
        http.Error(w, "Invalid dates", http.StatusBadRequest)
        return
    }

    // Hitung durasi menginap
    duration := int(checkOut.Sub(checkIn).Hours() / 24)

    // Mulai transaksi
    tx, err := db.Begin()
    if err != nil {
        http.Error(w, "Error starting transaction", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    var totalPrice float64
    var bookingIDs []int  // Pastikan ini slice []int

    // Proses setiap kamar yang dipesan
    for _, detail := range req.BookingDetails {
        log.Printf("Checking room_id: %d", detail.RoomID)

        // Ambil harga kamar berdasarkan tipe
        var pricePerNight float64
        query := `SELECT price_per_night FROM rooms WHERE room_id = ?`
        err := db.QueryRow(query, detail.RoomID).Scan(&pricePerNight)
        if err != nil {
            log.Printf("Error: Room ID '%d' not found. %v", detail.RoomID, err)
            http.Error(w, fmt.Sprintf("Error: Room ID '%d' not found. %v", detail.RoomID, err), http.StatusNotFound)
            return
        }

        // Hitung harga untuk tipe kamar ini
        totalRoomPrice := float64(detail.Quantity) * pricePerNight * float64(duration)
        totalPrice += totalRoomPrice

        // Simpan pemesanan untuk setiap tipe kamar
        query = `INSERT INTO bookings (user_id, room_id, check_in_date, check_out_date, total_price)
                 VALUES (?, ?, ?, ?, ?)`
        result, err := tx.Exec(query, req.CustomerID, detail.RoomID, req.CheckInDate, req.CheckOutDate, totalRoomPrice)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error booking room: %v", err), http.StatusInternalServerError)
            return
        }

        // Dapatkan ID pemesanan
        bookingID, err := result.LastInsertId()
        if err != nil {
            http.Error(w, "Error retrieving booking ID", http.StatusInternalServerError)
            return
        }
        bookingIDs = append(bookingIDs, int(bookingID))  // Menambahkan booking ID ke slice
    }

    // Simulasi pembayaran (misalnya menggunakan metode pembayaran tertentu)
    paymentMethod := req.PaymentDetails.PaymentMethod
    paymentStatus := "completed" // Status pembayaran sementara
    paymentDate := time.Now().Format("2006-01-02 15:04:05")

    // Simpan data pembayaran ke tabel payments
    for _, bookingID := range bookingIDs {
        query := `INSERT INTO payments (booking_id, payment_method, payment_status, payment_date, amount)
                  VALUES (?, ?, ?, ?, ?)`
        _, err := tx.Exec(query, bookingID, paymentMethod, paymentStatus, paymentDate, totalPrice)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error processing payment: %v", err), http.StatusInternalServerError)
            return
        }
    }

    // Commit transaksi
    err = tx.Commit()
    if err != nil {
        http.Error(w, "Error completing transaction", http.StatusInternalServerError)
        return
    }

    // Respons sukses
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(BookingResponse{
        BookingIDs: bookingIDs,  // Mengirim slice bookingIDs yang benar
        TotalPrice: totalPrice,
        Message:    "Rooms booked and payment processed successfully",
    })
}
