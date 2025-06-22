package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type SignInReq struct {
	Password string `json:"password"`
}

type SignInResp struct {
	Token string `json:"token"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var signInReq *SignInReq

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, &ErrorResp{"internal server error"})
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &signInReq)
	if err != nil {
		writeJson(w, &ErrorResp{"error unmarshal"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || signInReq.Password != pass {
		writeJson(w, &ErrorResp{"Неверный пароль"})
		return
	}

	secret := []byte(secretKey)

	claims := jwt.MapClaims{
		"password": createHash(pass),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		writeJson(w, &ErrorResp{"failed to sign jwt"})
		return
	}
	writeJson(w, &SignInResp{signedToken})
}
