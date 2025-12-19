-- +goose Up
CREATE TABLE favorites (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    item_type VARCHAR(50) NOT NULL,
    item_name VARCHAR(200) NOT NULL,
    is_favorite BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_favorites_user ON favorites(user_id);
CREATE INDEX idx_favorites_type ON favorites(item_type);

-- +goose Down
DROP TABLE favorites;
