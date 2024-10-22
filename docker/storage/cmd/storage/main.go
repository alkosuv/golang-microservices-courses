package main

import (
	"context"
	"flag"
	review "github.com/alkosuv/golang-microservices-courses/docker/storage"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/config"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/handler"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/migration"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/repository"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

// ./docker/storage/config/config.example.yml - для запуска из корня проекта
const (
	_configPath = "./config/config.example.yml"
)

func main() {
	ctx := context.Background()

	configFile := flag.String("config", _configPath, "config file path")
	flag.Parse()

	// Инициализация конфига
	cfg, err := config.NewConfig(*configFile)
	if err != nil {
		slog.Error("init config", "error", err)
		return
	}

	// Запуск миграции данных
	if err = migration.Up(ctx, cfg.DatabaseURL(), review.EmbedMigrations); err != nil {
		slog.Error("migrations up", "error", err)
		return
	}

	db, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	repo := repository.New(db)
	svc := service.New(repo)

	h := handler.New(svc)
	if err := h.ListenAndServe(cfg.App.HTTP.Host, cfg.App.HTTP.Port); err != nil {
		slog.Error("failed to start http server", "error", err)
	}
}
