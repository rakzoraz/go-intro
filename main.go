package main

import (
	"go-intro/internal/handler"
	"go-intro/internal/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	logger := logrus.NewEntry(logrus.New())
	h := &handler.Handler{}
	serv, err := server.New(8080, h)
	if err != nil {
		logger.WithField("error", err).Fatal("failed to init server")
	}

	serv.ListenAndServe(logger)

	// This allows us to listen for interrupts (ctrl+c, shutting down the run in goland/vscode, etc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("signal interrupt detected, shutting down ...")

	// shutdown the server
	if err := serv.Shutdown(); err != nil {
		logger.Fatalf("failed to shutdown server with error: %v", err)
	}

	// there was an interrupt so exit with code 1
	os.Exit(1)
}
