package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"auth-proxy/platform/authenticator"
	"auth-proxy/platform/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found. Please add a .env file with the required environment variables.")
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	log.Print("Server listening on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
