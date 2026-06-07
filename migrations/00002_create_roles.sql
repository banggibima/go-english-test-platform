-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO roles (id, name, description)
VALUES
  (gen_random_uuid(), 'admin', 'System Administrator'),
  (gen_random_uuid(), 'user', 'Regular User')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS roles;