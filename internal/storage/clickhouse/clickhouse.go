package clickhouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Константы для описания операций
const (
	opConn            = "clickhouse.Connect"
	opCloseConnection = "clickhouse.closeConnection"
	opParseDSN        = "clickhouse.ParseDSN"
	opPing            = "clickhouse.Ping"
	opInsert          = "clickhouse.Insert"
)

// Storage описывает слой взаимодействия с базой данных
type Storage struct {
	log *slog.Logger
	DB  driver.Conn
}

// NewStorage конструктор для Storage
func NewStorage(log *slog.Logger, dsn string) (*Storage, error) {
	// Парсим строку подключения
	options, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opParseDSN))
		return nil, fmt.Errorf("%s: %w", opParseDSN, err)
	}

	// Открываем соединение с базой данных
	db, err := clickhouse.Open(options)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opConn))
		return nil, fmt.Errorf("%s: %w", opConn, err)
	}

	// Проверяем подключение к базе данных
	if err = db.Ping(context.Background()); err != nil {
		log.Error("error", err.Error(), slog.String("operation", opPing))
		return nil, fmt.Errorf("%s: %w", opPing, err)
	}

	return &Storage{log: log, DB: db}, nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() {
	if err := s.DB.Close(); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opCloseConnection))
	}
}

// InsertRegistrationAnalytics сохраняет данные зарегистрировавшегося на мероприятие пользователя
func (s *Storage) InsertRegistrationAnalytics(ctx context.Context, chatID int64, username string, eventID string, createdAt time.Time) error {
	query := `INSERT INTO RegistrationAnalytics (chat_id, username, event_id, created_at) VALUES (?, ?, '?', ?)`
	dbCtx, _ := context.WithTimeout(ctx, 5*time.Second)
	batch, err := s.DB.PrepareBatch(dbCtx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", opInsert, err)
	}

	err = batch.Append(chatID, username, eventID, createdAt)
	if err != nil {
		return fmt.Errorf("%s: %w", opInsert, err)
	}

	if err = batch.Send(); err != nil {
		return fmt.Errorf("%s: %w", opInsert, err)
	}

	return nil
}
