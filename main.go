package main

import (
	"fmt"

	"github.com/jasrags/JasonMUD/server"
	"go.uber.org/zap"
)

const (
	ServerHost    = "localhost"
	ServerPort    = "3333"
	ServerNetwork = "tcp"
)

func main() {
	logger, errZap := zap.NewDevelopment()
	if errZap != nil {
		fmt.Printf("unable to setup logger: %v\n", errZap)
	}
	defer logger.Sync()

	s := server.New(logger)
	if err := s.Start(ServerHost, ServerPort, ServerNetwork); err != nil {
		logger.Panic("unable to start server", zap.Error(err))
	}
}
