package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("29bfe372865737fe2bfcfd3618b1da7d") // md5(mi)

/* -================================== Mock data ==================================-*/
var users = map[string]string{
	"catinho": "peixe",
	"lobinho": "ossinho",
}

/* -================================== Models ==================================-*/
// refac task: criar models
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

/* -================================== Handler  ==================================-*/
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
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
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
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
