package main

import (
	"log"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/commands"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}
	commands.Execute()
}
