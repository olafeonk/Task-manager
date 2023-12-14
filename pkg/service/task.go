package service

import (
	"task_manager"
	"task_manager/pkg/repository"
)

type TaskService struct {
	repo repository.TaskManagerTask
}

func NewTaskService(repo repository.TaskManagerTask) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(task task_manager.CreateTaskInput) (int, error) {
	return s.repo.Create(task, task_manager.Start)
}

func (s *TaskService) GetAll(telegramId int) ([]task_manager.Task, error) {
	return s.repo.GetAll(telegramId)
}

func (s *TaskService) GetById(taskId int) (task_manager.Task, error) {
	return s.repo.GetById(taskId)
}

func (s *TaskService) Delete(taskId int) error {
	return s.repo.Delete(taskId)
}
