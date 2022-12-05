package main

import (
	"go-kuebiko/kuebiko"
	"log"
	"os"
)

func main() {
	// Load environment variables
	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("Must set Discord token as env variable: BOT_TOKEN")
	}

	// Start the bot
	kuebiko.BotToken = botToken
	kuebiko.Run()

}
