package storage

import (
	"github.com/k1ender/task-master-go/internal/models"
	"gorm.io/gorm"
)

type TaskStore interface {
	CreateTask(task *models.Task) (*models.Task, error)
	GetTask(id uint) (*models.Task, error)
	GetTasks(userID uint) ([]models.Task, error)
	UpdateTask(destination *models.Task, updates map[string]any) error
	DeleteTask(id uint) error
}

type TaskStoreGorm struct {
	db *gorm.DB
}

func NewTaskStore(db *gorm.DB) TaskStore {
	return &TaskStoreGorm{db: db}
}

func (s *TaskStoreGorm) CreateTask(task *models.Task) (*models.Task, error) {
	return task, s.db.Create(task).Error
}

func (s *TaskStoreGorm) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	return &task, s.db.First(&task, id).Error
}

func (s *TaskStoreGorm) GetTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	return tasks, s.db.Where("user_id = ?", userID).Find(&tasks).Error
}

func (s *TaskStoreGorm) UpdateTask(destination *models.Task, updates map[string]any) error {
	return s.db.Model(destination).Updates(updates).Error
}

func (s *TaskStoreGorm) DeleteTask(id uint) error {
	return s.db.Delete(&models.Task{}, id).Error
}
