#!/bin/bash
set -e

echo "=== DNA Relations Complete Fix ==="
echo ""

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "Error: Run this script from /opt/dna-relations"
    exit 1
fi

echo "1. Running database migrations..."
docker exec -i dna-postgres psql -U dna -d dna << 'EOSQL'
-- Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL DEFAULT 'env_based',
    telegram_id BIGINT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);

-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) DEFAULT '#6366f1',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Topics
CREATE TABLE IF NOT EXISTS topics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    priority SMALLINT DEFAULT 3 CHECK (priority >= 0 AND priority <= 5),
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'discussed', 'archived')),
    created_by INT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    discussed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_topics_priority ON topics(priority);
CREATE INDEX IF NOT EXISTS idx_topics_status ON topics(status);
CREATE INDEX IF NOT EXISTS idx_topics_category ON topics(category_id);

-- Events
CREATE TABLE IF NOT EXISTS events (
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

CREATE INDEX IF NOT EXISTS idx_events_start ON events(start_at);
CREATE INDEX IF NOT EXISTS idx_events_owner ON events(owner_id);
CREATE INDEX IF NOT EXISTS idx_events_shared ON events(shared);

-- Notes
CREATE TABLE IF NOT EXISTS notes (
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

CREATE INDEX IF NOT EXISTS idx_notes_type ON notes(note_type);
CREATE INDEX IF NOT EXISTS idx_notes_pinned ON notes(is_pinned);

-- Favorites
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    item_type VARCHAR(50) NOT NULL,
    item_name VARCHAR(200) NOT NULL,
    is_favorite BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_favorites_user ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_type ON favorites(item_type);

-- Media sessions
CREATE TABLE IF NOT EXISTS media_sessions (
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

CREATE TABLE IF NOT EXISTS media_session_participants (
    session_id INT REFERENCES media_sessions(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (session_id, user_id)
);

-- Seed users
INSERT INTO users (id, username, password_hash) VALUES 
    (1, 'denis', 'env_based'),
    (2, 'nastya', 'env_based')
ON CONFLICT (id) DO NOTHING;

SELECT setval('users_id_seq', 2, true);

-- Seed categories
INSERT INTO categories (name, color) VALUES
    ('Отношения', '#ec4899'),
    ('Планы', '#6366f1'),
    ('Проблемы', '#ef4444'),
    ('Идеи', '#22c55e'),
    ('Разное', '#8b5cf6')
ON CONFLICT DO NOTHING;

SELECT 'Migrations complete!' as status;
EOSQL

echo ""
echo "2. Verifying tables..."
docker exec dna-postgres psql -U dna -d dna -c "\dt"

echo ""
echo "3. Restarting services..."
docker compose restart api bot

echo ""
echo "=== Done! ==="
echo ""
echo "Tables created. Now rebuild web with the fixed files:"
echo "  docker compose rm -f web"
echo "  docker compose up -d --build web"
echo ""
echo "Test the app at https://dna-relations.site"
