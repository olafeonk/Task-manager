package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task_manager"
	_ "task_manager/docs"
	"task_manager/pkg/handler"
	"task_manager/pkg/repository"
	"task_manager/pkg/service"
)

// @title Task Manager API
// @version 1.0
// @description API Server for TaskManager Application

// @host localhost:8080
// @BasePath /

func main() {
	//logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//slog.SetDefault(logger)
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableCaller = true
	zapConfig.Level.SetLevel(zap.DebugLevel)
	logger, _ := zapConfig.Build()

	http.Handle("/", handlerLog{logger})

	if err := godotenv.Load(); err != nil {
		slog.Error(fmt.Sprintf("error loading env variables: %s", err.Error()))
	}

	db, err := repository.NewPostgresDB(repository.Config{
		HOST:     os.Getenv("POSTGRES_HOST"),
		PORT:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(task_manager.Server)
	go func() {
		if err := server.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
			log.Panicf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Println("task manager started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("task manager shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error(fmt.Sprintf("error occured on server shutting down: %s", err.Error()))
	}

	if err := db.Close(); err != nil {
		slog.Error(fmt.Sprintf("error occured on db connection close: %s", err.Error()))
	}

}

type handlerLog struct {
	logger *zap.Logger
}

func (h handlerLog) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	h.logger.Info(
		"My log message",
		zap.String("my-key", "my-value"),
	)

	writer.WriteHeader(200)
	_, _ = writer.Write([]byte("Hello, world!"))
}
