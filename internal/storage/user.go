package storage

import (
	"github.com/k1ender/task-master-go/internal/models"
	"gorm.io/gorm"
)

type UserStore interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUser(id uint) (*models.User, error)
}

type UserStoreGorm struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &UserStoreGorm{db: db}
}

func (s *UserStoreGorm) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *UserStoreGorm) GetUser(id uint) (*models.User, error) {
	var user models.User
	return &user, s.db.First(&user, id).Error
}

func (s *UserStoreGorm) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	return &user, s.db.Where("username = ?", username).First(&user).Error
}
