package handlers

import (
	"cat-boxes-movies/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

/* -================================== Mock data ==================================-*/
var users = map[string]string{
	"catinho": "peixe",
	"lobinho": "ossinho",
}

/* -================================== Handler  ==================================-*/
func Login(w http.ResponseWriter, r *http.Request) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	fmt.Println(jwtKey)
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil || credentials.Username == "" || credentials.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// refac task: create postgress integration
	expectedPass, ok := users[credentials.Username]

	if !ok || expectedPass != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
	}
	// refac task: create function return clain
	expirationTime := time.Now().Add(time.Minute * 15)

	claims := &models.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
