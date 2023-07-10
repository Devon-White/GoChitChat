package main

import (
	"GoChitChat/server"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	server.Init()
}
