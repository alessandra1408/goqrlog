package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alessandra1408/goqrlog/app"
	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/alessandra1408/goqrlog/pkg/db"
	"github.com/alessandra1408/goqrlog/pkg/echoutil"
	"github.com/alessandra1408/goqrlog/pkg/httpclient"
	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/alessandra1408/goqrlog/routes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log := setupLogger()
	defer syncLogger(log)

	log.Info("Initializing API...")
	log.Info("Reading configurations...")
	cfg := setupConfig(log)

	log.Info("Configuring Echo framework...")
	e := echoutil.NewEcho()

	httpClient := httpclient.NewHTTPClient(cfg)

	log.Info("Setting up database connection...")

	db, err := db.NewDatabase(cfg.Database)
	if err != nil {
		log.Error("Error setting up database connection: ", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Conn.Close(context.Background()); err != nil {
			log.Warn("Error closing database connection: ", err)
		}
	}()

	log.Info("Configuring GO QR Logs Service...")

	appOpts := app.Options{
		HttpClient: httpClient,
		Cfg:        *cfg,
		DB:         db,
	}

	apps := app.New(appOpts)

	prefix := "/api"
	log.Infof("Configuring API routes in context_path: %s...", prefix)

	routesOptions := routes.Options{
		Group: e.Group(prefix),
		Apps:  apps,
		Cfg:   *cfg,
		Log:   log,
	}

	routes.RegisterRoutes(routesOptions)

	log.Info("Configuring server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	server := &http.Server{
		Addr:         ":" + port,
		IdleTimeout:  cfg.Server.IdleTimeout * time.Second,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
	}

	done := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Starting server...")
		err := e.StartServer(server)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Error on starting server: ", err)
		}
	}()

	<-done

	log.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := e.Shutdown(ctx); err != nil {
		log.Warn("Server shutdown failed: ", err)
	}

	log.Info("Server exited properly")
}

func setupLogger() log.Log {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := config.Build()
	if err != nil {
		fmt.Println("Unable to setup logger: " + err.Error())
		os.Exit(1)
	}

	sugar := logger.Sugar()

	return log.New(sugar)
}

func syncLogger(log log.Log) {
	if err := log.Sync(); err != nil {
		if !strings.Contains(err.Error(), "sync /dev/stdout: invalid argument") {
			log.Warn("Error syncing logger: " + err.Error())
		}
	}
}

func setupConfig(log log.Log) *config.Config {
	err := config.Log(log)
	if err != nil {
		log.Warn("Unable to setup logger level: ", err)
	}

	config.Local()

	cfg, err := config.Get()
	if err != nil {
		log.Error("Error on reading configurations: ", err)
		os.Exit(1)
	}

	return cfg
}
