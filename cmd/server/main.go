package main

import (
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/db"
	"github.com/k1ender/task-master-go/internal/models"
)

func main() {
	cfg := config.MustInit(".env")

	db := db.MustInit(cfg)
	db.AutoMigrate(&models.User{}, &models.Task{})
}
