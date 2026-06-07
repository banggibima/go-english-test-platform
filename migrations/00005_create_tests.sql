-- +goose Up
CREATE TABLE IF NOT EXISTS tests (
  id UUID PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  test_code VARCHAR(100) NOT NULL UNIQUE,
  duration_minutes INT NOT NULL,
  total_questions INT NOT NULL DEFAULT 0,
  passing_score INT NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_by UUID,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_tests_created_by
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON DELETE SET NULL
);

CREATE INDEX idx_tests_test_code ON tests(test_code);
CREATE INDEX idx_tests_is_active ON tests(is_active);
CREATE INDEX idx_tests_created_by ON tests(created_by);

-- +goose Down
DROP TABLE IF EXISTS tests;