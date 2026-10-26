// @title Coursepack API
// @version 1.0
// @description Curated lecture notes and past questions for university classes.
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/victorive/coursepack/docs"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/config"
	"github.com/victorive/coursepack/internal/server"
)

func main() {
	log := initLogger()

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic", "recover", r)
			os.Exit(1)
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	log.Info("starting", "app", cfg.App.Name, "env", cfg.App.Env)

	log.Info("initialising application container")
	container, err := app.New(context.Background(), cfg, log)
	if err != nil {
		log.Error("failed to initialise container", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := container.Close(); err != nil {
			log.Error("error closing container", "error", err)
		}
	}()

	httpServer := server.NewServer(cfg, container, log)
	httpServer.RegisterRoutes()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.Start(); err != nil {
			log.Error("server error", "error", err)
		}
	}()

	log.Info("server started", "port", cfg.App.Port, "swagger", "http://"+cfg.App.Host+":"+cfg.App.Port+"/docs/index.html")

	sig := <-sigChan
	log.Info("received signal", "signal", sig)

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(); err != nil {
		log.Error("server shutdown error", "error", err)
	}
	log.Info("server stopped gracefully")
}

func initLogger() *slog.Logger {
	level := slog.LevelInfo
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
