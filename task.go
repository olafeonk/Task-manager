package task_manager

import (
	"time"
)

type Task struct {
	Id          int        `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	TelegramId  int        `json:"telegram_id" db:"telegram_id"`
	Text        string     `json:"text" db:"text"`
	StartTimeAt time.Time  `json:"start_time_at" db:"start_time_at"`
	StatusEnd   StatusEnd  `json:"status_end" db:"status_end"`
	EndTask     *time.Time `json:"end_task_at" db:"end_task_at"`
}

type StatusEnd string

const (
	Start StatusEnd = "START"
	End   StatusEnd = "END"
)

type CreateTaskInput struct {
	Text         string `json:"text"`
	StartTime    time.Time
	StartTimeStr string `json:"start_time"`
	TelegramId   string `json:"telegram_id"`
}

type CreateTaskInputModeration struct {
	Text         string `form:"text" json:"text"`
	StartTimeStr []byte `form:"start_time" json:"start_time"`
	TelegramId   string `form:"telegram_id" json:"telegram_id"`
}
