package main

import (
	"cat-boxes-movies/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("working")
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/refresh", handlers.Refresh)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
