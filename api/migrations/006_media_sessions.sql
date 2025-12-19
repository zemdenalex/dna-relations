-- +goose Up
CREATE TABLE media_sessions (
    id SERIAL PRIMARY KEY,
    media_type VARCHAR(20) NOT NULL CHECK (media_type IN ('youtube', 'spotify', 'local')),
    media_url TEXT NOT NULL,
    media_title VARCHAR(300),
    host_user_id INT REFERENCES users(id),
    current_position_ms BIGINT DEFAULT 0,
    is_playing BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE media_session_participants (
    session_id INT REFERENCES media_sessions(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (session_id, user_id)
);

-- +goose Down
DROP TABLE media_session_participants;
DROP TABLE media_sessions;
