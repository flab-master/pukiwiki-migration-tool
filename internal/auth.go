package internal

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogin(adminUser, adminPass, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		userMatch := subtle.ConstantTimeCompare([]byte(req.Username), []byte(adminUser)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(req.Password), []byte(adminPass)) == 1
		if !userMatch || !passMatch {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		claims := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": signed})
	}
}

func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{},
				func(t *jwt.Token) (any, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, jwt.ErrSignatureInvalid
					}
					return []byte(jwtSecret), nil
				},
				jwt.WithExpirationRequired(),
			)
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
