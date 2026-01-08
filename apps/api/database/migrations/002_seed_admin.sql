-- +goose Up
-- Seed initial admin user (password: admin123 - change in production!)
-- Password hash is bcrypt of 'admin123'
INSERT INTO users (id, email, password_hash, role, display_name, is_active)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'admin@hive.local',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.uL1sKC8pPvV0Xq3wCi',
    'admin',
    'System Admin',
    true
);

-- +goose Down
DELETE FROM users WHERE email = 'admin@hive.local';
