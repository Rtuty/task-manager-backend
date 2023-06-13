package main

import (
	"log"
	"modules/internal/bot"
	"modules/internal/db"
	"modules/internal/todoist"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	go db.GetConnection()

	client, err := todoist.NewClient()
	if err != nil {
		panic(err)
	}

	bot.StartBot(client)
}
