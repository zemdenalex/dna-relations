-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) DEFAULT '#6366f1',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE topics (
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

CREATE TABLE topic_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE topic_tag_map (
    topic_id INT REFERENCES topics(id) ON DELETE CASCADE,
    tag_id INT REFERENCES topic_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (topic_id, tag_id)
);

CREATE INDEX idx_topics_priority ON topics(priority);
CREATE INDEX idx_topics_status ON topics(status);
CREATE INDEX idx_topics_category ON topics(category_id);

-- +goose Down
DROP TABLE topic_tag_map;
DROP TABLE topic_tags;
DROP TABLE topics;
DROP TABLE categories;
