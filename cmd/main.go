package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Astemirdum/logs/internal/config"
	"github.com/Astemirdum/logs/internal/handler"
	"github.com/Astemirdum/logs/internal/repository"
	"github.com/Astemirdum/logs/internal/service"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	configPath := pflag.StringP("config", "c", "config.yml", "config path")
	pflag.Parse()
	log := zap.NewExample()

	if err := godotenv.Load(); err != nil {
		log.Fatal("load envs from .env", zap.Error(err))
	}

	cfg := config.GetConfigYML(*configPath)

	db, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatal("db init: %v", zap.Error(err))
	}

	repo := repository.NewRepository(db, log)
	services := service.NewService(repo)
	h := handler.NewHandler(services, log)

	server := handler.NewServer(cfg.Server, h)

	server.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	termSig := <-sig

	log.Debug("Graceful shutdown starter", zap.Any("signal", termSig))

	closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_ = server.Stop(closeCtx)
	_ = db.Close()
}
