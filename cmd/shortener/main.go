package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RedWood011/ServiceURL/internal/config"
	"github.com/RedWood011/ServiceURL/internal/repository"
	"github.com/RedWood011/ServiceURL/internal/repository/memoryfile"
	"github.com/RedWood011/ServiceURL/internal/repository/postgres"
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp"
	"github.com/RedWood011/ServiceURL/internal/workers"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

// Global variables
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	var (
		db     *postgres.Repository
		dbFile *memoryfile.FileMap
		err    error
	)
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	cfg := config.NewConfig()
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr))

	if cfg.DatabaseDSN != "" {
		db, err = postgres.NewDatabase(ctx, cfg.DatabaseDSN, cfg.CountRepetitionBD)
	} else {
		dbFile, err = memoryfile.NewFileMap(cfg.FilePath)
	}
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(cfg.DatabaseDSN, db, dbFile)
	if err != nil {
		log.Fatal(err)
	}

	err = repo.Ping(ctx)
	if err != nil {
		log.Fatal("repo ping failed")
	}

	workerPool := workers.New(cfg.AmountWorkers, cfg.SizeBufWorker)

	serv := service.New(repo, logger, workerPool, cfg.BaseURL)

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS11,
		CurvePreferences: []tls.CurveID{},
		CipherSuites:     []uint16{},
	}

	httpServer := http.Server{
		Handler:   deliveryhttp.NewRouter(chi.NewRouter(), serv, cfg.KeyHash),
		Addr:      cfg.ServerAddress,
		TLSConfig: tlsConfig,
	}

	go func() {
		workerPool.Run(ctx)
	}()

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(ctx)

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

	if cfg.IsHTTPS {
		err = httpServer.ListenAndServeTLS("server_crt.crt", "server_key.key")
	} else {
		err = httpServer.ListenAndServe()
	}

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
