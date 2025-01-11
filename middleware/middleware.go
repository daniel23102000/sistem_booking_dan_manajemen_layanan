package middleware

import (
    "net/http"
    "strings"
    "fmt"
    "github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("your_secret_key")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

        next.ServeHTTP(w, r)
    }
}