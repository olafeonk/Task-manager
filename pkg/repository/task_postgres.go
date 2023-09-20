package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"strings"
	"task_manager"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) Create(userId int, task task_manager.CreateTaskInput, status task_manager.StatusEnd) (id int, err error) {
	query := fmt.Sprintf("INSERT INTO %s (name, user_id, status_end) VALUES ($1, $2, $3) RETURNING id", tasksTable)
	row := r.db.QueryRow(query, task.Name, userId, status)
	err = row.Scan(&id)
	return
}

func (r *TaskPostgres) GetAll(userId int) ([]task_manager.Task, error) {
	var tasks []task_manager.Task

	query := fmt.Sprintf("SELECT id, name, status_end, created_at, updated_at, end_task_at, user_id FROM %s WHERE user_id = $1", tasksTable)
	err := r.db.Select(&tasks, query, userId)

	return tasks, err
}

func (r *TaskPostgres) GetById(userId, taskId int) (task_manager.Task, error) {
	var task task_manager.Task

	query := fmt.Sprintf("SELECT id, name, status_end, created_at, updated_at, end_task_at, user_id FROM %s WHERE user_id = $1 AND id = $2", tasksTable)
	err := r.db.Get(&task, query, userId, taskId)

	return task, err
}

func (r *TaskPostgres) Delete(userId, taskId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND id = $2", tasksTable)
	_, err := r.db.Exec(query, userId, taskId)

	return err
}

func (r *TaskPostgres) Update(userId, taskId int, input task_manager.UpdateTaskInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, input.Name)
		argId++
	}

	if input.Status != "" {
		setValues = append(setValues, fmt.Sprintf("status_end=$%d", argId))
		args = append(args, input.Status)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d AND user_id=$%d",
		tasksTable, setQuery, argId, argId+1)
	args = append(args, taskId, userId)

	slog.Debug(fmt.Sprintf("updateQuery: %s", query))
	slog.Debug(fmt.Sprintf("args: %s", args))

	_, err := r.db.Exec(query, args...)
	return err

}
