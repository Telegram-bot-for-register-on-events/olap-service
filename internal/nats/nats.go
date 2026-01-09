package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Telegram-bot-for-register-on-events/olap-service/internal/domain/models"
	"github.com/nats-io/nats.go"
)

// Константы для описания операций
const (
	opConn            = "nats.Connect"
	opCreateJetStream = "nats.CreateJetStreamContext"
	opUnmarshal       = "nats.Unmarshal"
	opAck             = "nats.Ack"
	opSubscribe       = "nats.Subscribe"
	opUnsubscribe     = "nats.Unsubscribe"
	opInsert          = "nats.Insert"
)

// RegistrationAnalytics определяет метод для взаимодействия с сервисным слоем
type RegistrationAnalytics interface {
	InsertRegistrationAnalytics(ctx context.Context, chatID int64, username string, eventID string, createdAt time.Time) error
}

// Nats описывает брокер сообщений NATS
type Nats struct {
	log  *slog.Logger
	Conn *nats.Conn
	JS   nats.JetStreamContext
	sub  *nats.Subscription
	ra   RegistrationAnalytics
}

// NewNats конструктор для NATS
func NewNats(log *slog.Logger, url string, ra RegistrationAnalytics) (*Nats, error) {
	// Подключаемся к NATS по URL
	nc, err := nats.Connect(url)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opConn))
		return nil, fmt.Errorf("%s: %w", opConn, err)
	}

	log.Info("operation", opConn, slog.String("connected to", url))

	// Включаем JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opCreateJetStream))
		return nil, fmt.Errorf("%s: %w", opCreateJetStream, err)
	}
	log.Info("operation", opCreateJetStream, "successfully!")

	return &Nats{
		log:  log,
		Conn: nc,
		JS:   js,
		ra:   ra,
	}, nil
}

// RegistrationOnEvent метод для обработки сообщений из соответствующего топика
func (n *Nats) RegistrationOnEvent(ctx context.Context, topic string) error {
	// Подписываемся на топик
	sub, err := n.JS.Subscribe(topic, func(msg *nats.Msg) {
		var user models.User
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			n.log.Error("error", err.Error(), slog.String("operation", opUnmarshal))
		}
		// Передаём данные в сервисный слой
		err := n.ra.InsertRegistrationAnalytics(ctx, user.ChatID, user.Username, user.EventID, time.Now())
		if err != nil {
			n.log.Error("error", err.Error(), slog.String("operation", opInsert))
		}
		// Акаем данные
		if err = msg.Ack(); err != nil {
			n.log.Error("error", err.Error(), slog.String("operation", opAck))
		}
	}, nats.Durable("olap-service-consumer"), nats.ManualAck())
	n.sub = sub
	if err != nil {
		return fmt.Errorf("%s: %w", opSubscribe, err)
	}

	return nil
}

// Unsubscribe отписывается от топика
func (n *Nats) Unsubscribe() {
	err := n.sub.Unsubscribe()
	if err != nil {
		n.log.Error("error", err.Error(), slog.String("operation", opUnsubscribe))
	}
}
