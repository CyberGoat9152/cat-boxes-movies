package main

import (
	"cat-boxes-movies/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	fmt.Println("working")
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/refresh", handlers.Refresh)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("API_PORT"), nil))
}
