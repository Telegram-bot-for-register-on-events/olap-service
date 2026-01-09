-- +goose Up
CREATE TABLE IF NOT EXISTS RegistrationAnalytics
(
    id UUID default generateUUIDv4(),
    chat_id UInt64,
    username Nullable(String),
    event_id UUID,
    created_at DATETIME64
)
    ENGINE = MergeTree()
        ORDER BY (event_id);

-- +goose Down
DROP TABLE IF EXISTS RegistrationAnalytics;
