package task_manager

import (
	"errors"
	"time"
)

type Task struct {
	Id        int        `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	UserId    int        `json:"user_id" db:"user_id"`
	Name      string     `json:"name" db:"name"`
	StatusEnd StatusEnd  `json:"status_end" db:"status_end"`
	EndTask   *time.Time `json:"end_task_at" db:"end_task_at"`
}

type StatusEnd string

const (
	Start StatusEnd = "START"
	End   StatusEnd = "END"
)

type UpdateTaskInput struct {
	Name   string    `json:"name"`
	Status StatusEnd `json:"status"`
}

type CreateTaskInput struct {
	Name string `json:"name" binding:"required"`
}

func (i UpdateTaskInput) Validate() error {
	if i.Name == "" && i.Status == "" {
		return errors.New("body has no values")
	}
	if i.Status != "" && i.Status != Start && i.Status != End {
		return errors.New("invalid status")
	}
	return nil
}
