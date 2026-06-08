-- +goose Up
CREATE TABLE IF NOT EXISTS files (
  id UUID PRIMARY KEY,
  user_id UUID,
  bucket VARCHAR(255) NOT NULL,
  object_name TEXT NOT NULL,
  original_name VARCHAR(255) NOT NULL,
  content_type VARCHAR(100) NOT NULL,
  size_bytes BIGINT NOT NULL,
  file_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_files_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE SET NULL
);

CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_object_name ON files(object_name);

-- +goose Down
DROP TABLE IF EXISTS files;