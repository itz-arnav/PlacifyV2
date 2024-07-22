package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func ValidateToken(tokenString string) (*Claims, bool) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, false
	}

	return claims, true
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Missing auth token"))
			return
		}

		claims, valid := ValidateToken(tokenHeader)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid or expired auth token"))
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
