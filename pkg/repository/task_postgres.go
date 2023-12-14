package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"task_manager"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) Create(task task_manager.CreateTaskInput, status task_manager.StatusEnd) (id int, err error) {
	query := fmt.Sprintf("INSERT INTO %s (text, telegram_id, status_end, start_time_at) VALUES ($1, $2, $3, $4) RETURNING id", tasksTable)
	row := r.db.QueryRow(query, task.Text, task.TelegramId, status, task.StartTime)
	err = row.Scan(&id)
	return
}

func (r *TaskPostgres) GetAll(telegramId int) ([]task_manager.Task, error) {
	var tasks []task_manager.Task

	query := fmt.Sprintf("SELECT id, text, status_end, created_at, updated_at, end_task_at, telegram_id, start_time_at FROM %s WHERE telegram_id = $1", tasksTable)
	err := r.db.Select(&tasks, query, telegramId)

	return tasks, err
}

func (r *TaskPostgres) GetById(taskId int) (task_manager.Task, error) {
	var task task_manager.Task

	query := fmt.Sprintf("SELECT id, text, status_end, created_at, updated_at, end_task_at, telegram_id, start_time_at FROM %s WHERE id = $1", tasksTable)
	err := r.db.Get(&task, query, taskId)

	return task, err
}

func (r *TaskPostgres) Delete(taskId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tasksTable)
	_, err := r.db.Exec(query, taskId)

	return err
}
