package models

// User описывает сообщения из NATS
type User struct {
	ChatID   int64  `json:"chat_id"`
	Username string `json:"username"`
	EventID  string `json:"event_id"`
}
