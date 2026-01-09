package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// Ещё одна точка входа, для выполнения миграций к базе данных
func main() {
	// Достаём из переменных окружения путь к миграциям и строку подключения
	var (
		driverName = "clickhouse"
		dsn        = os.Getenv("DSN")
		dir        = os.Getenv("DIR")
	)

	// Парсим флаги, определяем команду для goose
	flag.Parse()

	log := setupLogger()

	// Выполняем подключение к базе данных
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Error("error opening database connection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Пингуем базу данных
	if err = db.PingContext(context.Background()); err != nil {
		log.Error("error pinging database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Устанавливаем диалект для гуся
	_ = goose.SetDialect(driverName)

	args := flag.Args()
	if len(args) == 0 {
		log.Error("command is required: up, down, status")
		os.Exit(1)
	}

	command := args[0]

	// Запускаем миграции
	if err = goose.RunContext(context.Background(), command, db, dir, args[1:]...); err != nil {
		log.Error("error running migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	_ = db.Close()

	log.Info("migrations complete", slog.String("command", command))
}

// setupLogger инициализирует логгер с JSON-обработчиком
func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
