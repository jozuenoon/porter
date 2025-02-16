package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"porter"
	"porter/adapter"
	"porter/app/getport"
	"porter/app/portsfromstream"
	"porter/controller"
)

func main() {
	repo := adapter.NewInMemoryRepository()
	ingestor := adapter.NewIngestor[porter.Port]()
	batchSize := 1000 // to be configured via environment variable.

	portsFromStreamUsecase := portsfromstream.NewService(repo, ingestor, batchSize)
	getPortUsecase := getport.NewService(repo)

	svc := &struct {
		*portsfromstream.PortsFromStreamService
		*getport.GetPortService
	}{
		portsFromStreamUsecase,
		getPortUsecase,
	}

	httpController := controller.NewHTTPController(svc)

	server := &http.Server{Addr: ":8080", Handler: httpController}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server died.")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	const graceSeconds time.Duration = 5

	ctx, cancel := context.WithTimeout(context.Background(), graceSeconds*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown.")
	}

	log.Info().Msg("Server gracefully stopped")
}
