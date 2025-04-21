package storage

import "gorm.io/gorm"

type Storage struct {
	Users UserStore
	Tasks TaskStore
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		Users: NewUserStore(db),
		Tasks: NewTaskStore(db),
	}
}
