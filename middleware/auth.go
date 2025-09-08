// File: middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/keshav7976/ecommerce/utils"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Missing token"}`, http.StatusUnauthorized)
			return
		}

		// Split the header to handle different cases
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := utils.ParseJWT(token)
		if err != nil {
			http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, `{"error":"Invalid token claims"}`, http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), "user_id", uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}