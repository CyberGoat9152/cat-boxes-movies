package handlers

import (
	"cat-boxes-movies/models"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

/* -================================== Mock data ==================================-*/
var users = map[string]string{
	"catinho": "peixe",
	"lobinho": "ossinho",
}

func Login(w http.ResponseWriter, r *http.Request) {
	// Variables
	var credentials models.Credentials
	var tokenString = ""
	expirationTime := time.Now()

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
	status := generateToken(credentials, &tokenString, &expirationTime)
	// Cookie setting to response
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	// Header and response configuration
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func generateToken(credentials models.Credentials, tokenStr *string, expiration *time.Time) int {
	// Calc expiration time
	jwtExpirationTime, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME"))
	expirationTime := time.Now().Add(time.Minute * time.Duration(jwtExpirationTime))
	// Create claim var
	claims := &models.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	// Pass values
	*tokenStr = tokenString
	*expiration = expirationTime
	// Error Handling
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
