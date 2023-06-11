package main

import (
	"log"
	"modules/internal/db"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	db.GetConnection()
}
