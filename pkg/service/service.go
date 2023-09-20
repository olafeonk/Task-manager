package service

import (
	"task_manager"
	"task_manager/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user task_manager.User) (int, error)
	GetUserId(username, password string) (int, error)
}

type TaskManagerTask interface {
	Create(userId int, task task_manager.CreateTaskInput) (int, error)
	GetAll(userId int) ([]task_manager.Task, error)
	GetById(userId, taskId int) (task_manager.Task, error)
	Delete(userId, taskId int) error
	Update(userId, taskId int, input task_manager.UpdateTaskInput) error
}

type Service struct {
	Authorization
	TaskManagerTask
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:   NewAuthService(repos.Authorization),
		TaskManagerTask: NewTaskService(repos.TaskManagerTask),
	}
}
