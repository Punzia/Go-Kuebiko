package main

import (
	"go-kuebiko/kuebiko"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get the bot token from .env! Space for API tokens...
	botToken := os.Getenv("BOT_TOKEN")

	// Start the bot
	kuebiko.BotToken = botToken
	kuebiko.Run()

}
