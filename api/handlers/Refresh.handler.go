package handlers

import (
	"cat-boxes-movies/models"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Refresh(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")

	if err != nil {
		// cookie exist ?
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// else bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// it's a valid token ?
	tokenStr := cookie.Value
	claims := &models.Claims{}
	expirationTime := time.Now()
	status := 500
	if validateToken(&tokenStr, claims) {
		status = refreshToken(&tokenStr, &expirationTime, claims)
	}
	//with not expired
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 45*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Cookie setting to response
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	})
	// Header and response configuration
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenStr,
	})
}

func validateToken(tokenStr *string, claims *models.Claims) bool {
	tkn, _ := jwt.ParseWithClaims(*tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return tkn.Valid
}

func refreshToken(tokenStr *string, expiration *time.Time, claims *models.Claims) int {
	jwtExpirationTime, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME"))
	expirationTime := time.Now().Add(time.Minute * time.Duration(jwtExpirationTime))
	// Refresh the expire time
	claims.ExpiresAt = expirationTime.Unix()
	// Generate new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	// Fill variables
	*tokenStr = tokenString
	*expiration = expirationTime
	// Handle errors
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
