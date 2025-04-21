package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	Title     string       `json:"title" gorm:"not null"`
	Body      string       `json:"body" gorm:"not null"`
	Completed sql.NullBool `json:"completed" gorm:"default:false"`
	UserID    uint         `json:"-" gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
