-- +goose Up
CREATE TABLE IF NOT EXISTS RegistrationAnalytics
(
    id UUID default generateUUIDv4(),
    chat_id Int64,
    username String,
    event_id String,
    created_at DATETIME64
)
    ENGINE = MergeTree()
        ORDER BY (event_id);

-- +goose Down
DROP TABLE IF EXISTS RegistrationAnalytics;
