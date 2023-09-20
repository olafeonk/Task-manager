package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"task_manager"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user task_manager.User) (id int, err error) {
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Username, user.Password)
	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, err
}

func (r *AuthPostgres) GetUser(username, password string) (user task_manager.User, err error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err = r.db.Get(&user, query, username, password)
	return
}
