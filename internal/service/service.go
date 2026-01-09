package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// Константы для описания операций
const (
	opInsertRA = "service.InsertRA"
)

// Service описывает сервисный слой микросервиса
type Service struct {
	log *slog.Logger
	ra  RegistrationAnalytics
}

// RegistrationAnalytics описывает метод для взаимодействия с базой данных
type RegistrationAnalytics interface {
	InsertRegistrationAnalytics(ctx context.Context, chatID int64, username string, eventID string, createdAt time.Time) error
}

// NewService конструктор для создания Service
func NewService(log *slog.Logger, ra RegistrationAnalytics) *Service {
	return &Service{
		log: log,
		ra:  ra,
	}
}

// InsertRegistrationAnalytics валидирует входящие данные и передаёт их в слой взаимодействия с базой данных
func (s *Service) InsertRegistrationAnalytics(ctx context.Context, chatID int64, username string, eventID string, createdAt time.Time) error {
	if err := validateChatID(chatID); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opInsertRA))
		return err
	}

	if err := validateUsername(username); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opInsertRA))
		return err
	}

	if err := validateEventID(eventID); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opInsertRA))
		return err
	}

	err := s.ra.InsertRegistrationAnalytics(ctx, chatID, username, eventID, createdAt)
	if err != nil {
		return fmt.Errorf("%s: %w", opInsertRA, err)
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	} else if len(username) < 5 || len(username) > 32 {
		return errors.New("username length must be between 5 and 32")
	}
	return nil
}

func validateChatID(chatID int64) error {
	if chatID == 0 {
		return errors.New("chatID cannot be equal to 0")
	} else if chatID < -999999999999999 || chatID > 999999999999999 {
		return errors.New("chatID out of range")
	}
	return nil
}

func validateEventID(eventID string) error {
	if eventID == "" {
		return errors.New("eventID cannot be empty")
	}
	return nil
}
