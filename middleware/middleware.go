package middleware

import (
    "database/sql"
    "net/http"
    "strings"
    "fmt"
    "github.com/golang-jwt/jwt/v4"
    _ "github.com/go-sql-driver/mysql" // Pastikan driver MySQL sudah terpasang
)

var secretKey = []byte("your_secret_key")

// Fungsi untuk mendapatkan email dari klaim token JWT
func getEmailFromClaims(claims jwt.MapClaims) (string, error) {
    email, ok := claims["email"].(string)
    if !ok {
        return "", fmt.Errorf("email not found in token")
    }
    return email, nil
}

// Fungsi untuk mendapatkan role pengguna berdasarkan email dari database
func getRoleFromEmail(db *sql.DB, email string) (string, error) {
    var role string
    query := `SELECT role FROM users WHERE email = ?`
    err := db.QueryRow(query, email).Scan(&role)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("user not found")
        }
        return "", fmt.Errorf("error fetching user role from database: %v", err)
    }
    return role, nil
}

// AuthMiddleware untuk melindungi route berdasarkan role pengguna
func AuthMiddleware(allowedRoles []string, db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return secretKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, fmt.Sprintf("Invalid Token: %v", err), http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid claims", http.StatusUnauthorized)
            return
        }

        // Menggunakan fungsi untuk mengambil email dari klaim
        email, err := getEmailFromClaims(claims)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        // Menggunakan fungsi untuk mendapatkan role berdasarkan email
        role, err := getRoleFromEmail(db, email)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        // Verifikasi apakah role sesuai dengan role yang diizinkan
        if !contains(allowedRoles, role) {
            http.Error(w, "Forbidden: Insufficient role", http.StatusForbidden)
            return
        }

        // Menambahkan email ke dalam header respons jika diperlukan
        w.Header().Set("X-User-Email", email)

        // Melanjutkan eksekusi request
        next.ServeHTTP(w, r)
    }
}

// Helper function untuk memeriksa apakah role ada dalam allowedRoles
func contains(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}