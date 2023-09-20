package repository

import (
	"github.com/jmoiron/sqlx"
	"task_manager"
)

type Authorization interface {
	CreateUser(user task_manager.User) (int, error)
	GetUser(username, password string) (task_manager.User, error)
}

type TaskManagerTask interface {
	Create(userId int, task task_manager.CreateTaskInput, status task_manager.StatusEnd) (int, error)
	GetAll(userId int) ([]task_manager.Task, error)
	GetById(userId, taskId int) (task_manager.Task, error)
	Delete(userId, taskId int) error
	Update(userId, taskId int, input task_manager.UpdateTaskInput) error
}

type Repository struct {
	Authorization
	TaskManagerTask
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:   NewAuthPostgres(db),
		TaskManagerTask: NewTaskPostgres(db),
	}
}
