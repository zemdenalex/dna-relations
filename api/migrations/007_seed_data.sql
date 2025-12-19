-- +goose Up
INSERT INTO users (id, username, password_hash) VALUES 
    (1, 'denis', 'env_based'),
    (2, 'nastya', 'env_based')
ON CONFLICT (id) DO NOTHING;

SELECT setval('users_id_seq', 2);

INSERT INTO categories (name, color) VALUES
    ('Relationship', '#ec4899'),
    ('Plans', '#6366f1'),
    ('Problems', '#ef4444'),
    ('Ideas', '#22c55e'),
    ('Random', '#8b5cf6')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM categories;
DELETE FROM users WHERE id IN (1, 2);
