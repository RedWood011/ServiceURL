package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RedWood011/ServiceURL/internal/config"
	"github.com/RedWood011/ServiceURL/internal/repository"
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.NewConfig()
	logger := slog.New(slog.NewTextHandler(os.Stderr))
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}

	serv := service.New(repo, logger, cfg.Address)

	httpServer := http.Server{
		Handler: deliveryhttp.NewRouter(chi.NewRouter(), serv, cfg.KeyHash),
		Addr:    cfg.Port,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)

		err = serv.Repo.SaveDone()
		if err != nil {
			log.Fatal("SaveDone error")
		}

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err = httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
		cancel()
	}()

	if err = httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
