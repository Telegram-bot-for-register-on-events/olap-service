package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	opLoadConfig = "config.load"
)

// Config описывает конфигурацию микросервиса
type Config struct {
	natsCfg *natsConfig
	dbCfg   *databaseConfig
}

// natsConfig описывает конфигурацию брокера сообщений NATS
type natsConfig struct {
	url   string
	topic string
}

// databaseConfig описывает конфигурацию базы данных
type databaseConfig struct {
	path string
}

// getEnv проверяет наличие переменной окружения и возвращает её текущее значение, либо стандартное, при отсутствии текущего
func getEnv(key, reserve string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return reserve
}

// newNatsConfig загружает конфигурацию для NATS
func newNatsConfig(log *slog.Logger) (*natsConfig, error) {
	url := getEnv("NATS_URL", "")
	if url == "" {
		log.Error("nats url cannot be empty")
		return nil, errors.New("nats url cannot be empty")
	}

	topic := getEnv("NATS_TOPIC", "")
	if topic == "" {
		log.Error("nats topic cannot be empty")
		return nil, errors.New("nats topic cannot be empty")
	}

	return &natsConfig{url: url, topic: topic}, nil
}

// newDatabaseConfig загружает конфигурацию для базы данных
func newDatabaseConfig(log *slog.Logger) (*databaseConfig, error) {
	path := getEnv("DSN", "")
	if path == "" {
		log.Error("dsn cannot be empty")
		return nil, errors.New("dsn cannot be empty")
	}

	return &databaseConfig{path: path}, nil
}

// LoadConfig конструктор для Config
func LoadConfig(log *slog.Logger) (*Config, error) {
	log.Info("loading environment variables")
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Error("operation", opLoadConfig, err.Error())
		return nil, fmt.Errorf("%s: %w", opLoadConfig, err)
	}
	log.Info("environment variables successfully loaded")

	// Создаём конфигурацию базы данных
	dbCfg, err := newDatabaseConfig(log)
	if err != nil {
		log.Error("operation", opLoadConfig, err.Error())
		return nil, err
	}
	// Создаём конфигурацию NATS
	natsCfg, err := newNatsConfig(log)
	if err != nil {
		log.Error("operation", opLoadConfig, err.Error())
		return nil, err
	}

	return &Config{natsCfg: natsCfg, dbCfg: dbCfg}, nil
}

// MustLoadConfig обёртка для LoadConfig - при ошибке - паникует
func MustLoadConfig(log *slog.Logger) *Config {
	cfg, err := LoadConfig(log)
	if err != nil {
		panic(err)
	}
	return cfg
}

// GetNatsURL геттер для получения URL для подключения к NATS
func (c *Config) GetNatsURL() string {
	return c.natsCfg.url
}

// GetDatabasePath геттер, для получения пути подключения к базе данных
func (c *Config) GetDatabasePath() string {
	return c.dbCfg.path
}

// GetNatsTopic геттер для получения названия топика, в который будут публиковаться сообщения
func (c *Config) GetNatsTopic() string { return c.natsCfg.topic }
