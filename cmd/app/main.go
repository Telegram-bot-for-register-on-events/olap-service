package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Telegram-bot-for-register-on-events/olap-service/internal/app"
)

func main() {
	// Инициализируем логгер
	log := setupLogger()
	// Создаём новый инстанс микросервиса
	application := app.NewApp(log)
	// Запускаем его
	application.MustStart()
	// Создаём канал для приёма сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	// Передаём входящие сигналы в канал stop
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// Читаем из канала, пока не придёт соответствующий сигнал
	<-sigChan

	application.Stop()
}

// setupLogger инициализирует логгер с JSON-обработчиком
func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
