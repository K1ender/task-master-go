package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"not null"`
	Tasks     []Task `json:"tasks" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
