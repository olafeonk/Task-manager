package repository

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log/slog"
	"os"
)

const (
	tasksTable = "tasks"
)

type Config struct {
	HOST     string
	PORT     string
	Username string
	Password string
	DBName   string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	slog.Info("init database",
		"cfg", cfg)
	rootCertPool := x509.NewCertPool()
	pathSert := "./root.crt"
	slog.Info(pathSert)
	slog.Info("перед чтением файла")
	pem, err := ioutil.ReadFile(pathSert)
	if err != nil {
		panic(err)
	}
	slog.Info("после чтения файла")
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		panic("Failed to append PEM.")
	}
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=verify-full sslrootcert=%s target_session_attrs=read-write",
		cfg.HOST, cfg.PORT, cfg.Username, cfg.DBName, cfg.Password, pathSert)
	slog.Info(connString)
	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}
	connConfig.TLSConfig = &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	var version string

	err = conn.QueryRow(context.Background(), "select version()").Scan(&version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(version)

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://schema", "postgres", driver)
	if err != nil {
		return nil, err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return db, nil
}
