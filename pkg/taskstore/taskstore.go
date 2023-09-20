package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
type Task struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserId    int       `json:"userId"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	EndTask   time.Time `json:"endTask"`
}

type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts
}

// CreateTask creates a new task in the store.
func (ts *TaskStore) CreateTask(text string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id:      ts.nextId,
		Name:    text,
		EndTask: due,
	}

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

// GetTask retrieves a task from the store, by id. If no such id exists, an
// error is returned.
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

// DeleteTask deletes the task with the given id. If no such id exists, an error
// is returned.
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}
