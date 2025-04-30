package main

import (
	"net/http"

	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/db"
	"github.com/k1ender/task-master-go/internal/logger"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/routes"
	"github.com/k1ender/task-master-go/internal/storage"
)

// @title Task Master API
// @description Task Master API - Simple task manager
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustInit(".env")

	db := db.MustInit(cfg)
	db.AutoMigrate(&models.User{}, &models.Task{})

	storage := storage.NewStorage(db)

	logger := logger.MustInit(cfg)

	router := routes.New(db, cfg, storage, logger)

	logger.Info("Server started", "port", cfg.HttpServer.Port)

	http.ListenAndServe(":"+cfg.HttpServer.Port, router)
}
