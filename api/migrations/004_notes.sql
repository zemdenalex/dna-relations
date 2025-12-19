-- +goose Up
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200),
    content TEXT NOT NULL,
    note_type VARCHAR(30) DEFAULT 'general' CHECK (note_type IN ('general', 'rule', 'thought', 'resource', 'credential')),
    is_pinned BOOLEAN DEFAULT FALSE,
    created_by INT REFERENCES users(id),
    shared BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notes_type ON notes(note_type);
CREATE INDEX idx_notes_pinned ON notes(is_pinned);

-- +goose Down
DROP TABLE notes;
