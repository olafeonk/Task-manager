package repository

import (
	"github.com/jmoiron/sqlx"
	"task_manager"
)

type TaskManagerTask interface {
	Create(task task_manager.CreateTaskInput, status task_manager.StatusEnd) (int, error)
	GetAll(telegramId int) ([]task_manager.Task, error)
	GetById(taskId int) (task_manager.Task, error)
	Delete(taskId int) error
}

type Repository struct {
	TaskManagerTask
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TaskManagerTask: NewTaskPostgres(db),
	}
}
