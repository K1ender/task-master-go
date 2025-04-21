package db

import (
	"fmt"

	"github.com/k1ender/task-master-go/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustInit(config *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(DSN(config)), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func DSN(config *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Name)
}
