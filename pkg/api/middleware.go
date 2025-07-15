package api

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(pass) > 0 {
			var token string
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
			var valid bool
			secret := []byte(secretKey)

			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				fmt.Printf("failed to parse token: %s\n", err)
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			res, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				fmt.Printf("failed to typecast to jwt.MapClaims")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			passwordHash := res["password"]
			valid = passwordHash == createHash(pass)
			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
