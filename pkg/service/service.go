package service

import (
	"task_manager"
	"task_manager/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type TaskManagerTask interface {
	Create(task task_manager.CreateTaskInput) (int, error)
	GetAll(telegramId int) ([]task_manager.Task, error)
	GetById(taskId int) (task_manager.Task, error)
	Delete(taskId int) error
}

type Service struct {
	TaskManagerTask
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		TaskManagerTask: NewTaskService(repos.TaskManagerTask),
	}
}
