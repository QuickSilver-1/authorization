package main

import (
	"auth/internal/config"
	"auth/internal/logger"
	"auth/internal/server"
)

func main() {
	err := config.NewConfig()

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	s := server.NewServer(config.AppConfig.HttpPort, 10, 10)

	err = s.StartServer()

	if err != nil {
		logger.Log.Error(err.Error())
	}
}