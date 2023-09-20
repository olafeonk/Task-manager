package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"task_manager"
	"task_manager/pkg/repository"
	"task_manager/pkg/service"
)

func main() {
	var (
		username string
		password string
	)
	if err := godotenv.Load(); err != nil {
		slog.Error(fmt.Sprintf("error loading env variables: %s", err.Error()))
	}
	cfg := repository.Config{
		HOST:     os.Getenv("POSTGRES_HOST"),
		PORT:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password user")
	flag.Parse()
	slog.Debug("Database",
		"config", cfg)
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.HOST, cfg.PORT, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize db: %s", err.Error()))
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("failed to ping db: %s", err.Error()))
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	user := task_manager.User{
		Username: username,
		Password: password,
	}
	id, err := services.Authorization.CreateUser(user)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create user: %s", err.Error()))
	}
	slog.Info("success",
		"id", id)
}
