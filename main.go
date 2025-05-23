package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/berryscottr/home-assistant/api"
	"github.com/berryscottr/home-assistant/devices"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Starting Local Home Automation System")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devices.Init(ctx)
	go api.StartServer(ctx)

	// Wait for termination signal (Ctrl+C or kill)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info().Msg("Shutting down")
}
