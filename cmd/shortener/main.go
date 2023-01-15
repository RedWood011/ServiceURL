package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RedWood011/ServiceURL/internal/repository"
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp"
	"github.com/go-chi/chi/v5"
)

const address = "http://localhost:8080"

func main() {
	repo := repository.NewRepository("")

	serv := service.New(repo, address)
	handler := chi.NewRouter()
	deliveryhttp.NewRouter(handler, serv)
	httpServer := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 5 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 5*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
