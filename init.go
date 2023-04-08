package main

import (
	"github.com/PretendoNetwork/plogger-go"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/joho/godotenv"
)

var logger = plogger.NewLogger()

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Warning("Error loading .env file")
	}

	database.ConnectAll()
}
