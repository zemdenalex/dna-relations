-- +goose Up
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ,
    all_day BOOLEAN DEFAULT FALSE,
    owner_id INT REFERENCES users(id),
    shared BOOLEAN DEFAULT TRUE,
    color VARCHAR(7) DEFAULT '#6366f1',
    recurrence VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_events_start ON events(start_at);
CREATE INDEX idx_events_owner ON events(owner_id);
CREATE INDEX idx_events_shared ON events(shared);

-- +goose Down
DROP TABLE events;
