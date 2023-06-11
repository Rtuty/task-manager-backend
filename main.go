package main

import (
	"context"
	"log"
	"modules/internal/bot"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot.StartBotInstance(ctx)
}
