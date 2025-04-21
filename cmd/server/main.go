package main

import (
	"net/http"

	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/db"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/routes"
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

	router := routes.New(db, cfg)

	http.ListenAndServe(":"+cfg.HttpServer.Port, router)
}
