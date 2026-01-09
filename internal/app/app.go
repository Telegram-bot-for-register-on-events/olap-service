package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/Telegram-bot-for-register-on-events/olap-service/internal/config"
	ns "github.com/Telegram-bot-for-register-on-events/olap-service/internal/nats"
	"github.com/Telegram-bot-for-register-on-events/olap-service/internal/service"
	"github.com/Telegram-bot-for-register-on-events/olap-service/internal/storage/clickhouse"
)

// App описывает микросервис целиком
type App struct {
	log  *slog.Logger
	Nats *ns.Nats
	DB   *clickhouse.Storage
	Cfg  *config.Config
}

// NewApp конструктор для App
func NewApp(log *slog.Logger) *App {
	// Загружаем конфигурацию
	cfg := cfgInit(log)
	// Инициализируем хранилище данных
	db := dbInit(log, cfg.GetDatabasePath())
	// Инициализируем сервисный слой
	srvc := service.NewService(log, db)
	// Подключаемся к Nats
	nats := natsConn(log, cfg.GetNatsURL(), srvc)
	return &App{log: log, Nats: nats, DB: db, Cfg: cfg}
}

// MustStart подписывается на топик и начинает обработку сообщений из него, при ошибке - паникует
func (a *App) MustStart() {
	err := a.Nats.RegistrationOnEvent(context.Background(), a.Cfg.GetNatsTopic())
	if err != nil {
		panic(err)
	}
}

// Stop выполняется при Graceful Shutdown
func (a *App) Stop() {
	a.log.Info("shutting down...")
	a.Nats.Unsubscribe()
	a.Nats.Conn.Close()
	a.DB.Close()
}

// cfgInit обёртка для инициализации конфига
func cfgInit(log *slog.Logger) *config.Config {
	cfg := config.MustLoadConfig(log)
	log.Info("config successfully loaded")
	return cfg
}

// dbInit обёртка для создания подключения к базе данных
func dbInit(log *slog.Logger, dsn string) *clickhouse.Storage {
	db, err := clickhouse.NewStorage(log, dsn)
	if err != nil {
		log.Error("error", err.Error(), slog.String("failed", "connect to database"))
		os.Exit(1)
	}
	log.Info("connection to database successfully")
	return db
}

// natsConn обёртка для подключения к NATS
func natsConn(log *slog.Logger, url string, ra ns.RegistrationAnalytics) *ns.Nats {
	nats, err := ns.NewNats(log, url, ra)
	if err != nil {
		os.Exit(1)
	}
	log.Info("connected to nats successfully")
	return nats
}
