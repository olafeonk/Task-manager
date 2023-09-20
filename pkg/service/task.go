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

func (s *TaskService) Create(userId int, task task_manager.CreateTaskInput) (int, error) {
	return s.repo.Create(userId, task, task_manager.Start)
}

func (s *TaskService) GetAll(userId int) ([]task_manager.Task, error) {
	return s.repo.GetAll(userId)
}

func (s *TaskService) GetById(userId, taskId int) (task_manager.Task, error) {
	return s.repo.GetById(userId, taskId)
}

func (s *TaskService) Delete(userId, taskId int) error {
	return s.repo.Delete(userId, taskId)
}

func (s *TaskService) Update(userId, taskId int, input task_manager.UpdateTaskInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, taskId, input)

}
