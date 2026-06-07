-- +goose Up
CREATE TABLE IF NOT EXISTS results (
  id UUID PRIMARY KEY,
  attempt_id UUID NOT NULL UNIQUE,
  user_id UUID NOT NULL,
  test_id UUID NOT NULL,
  score INT NOT NULL DEFAULT 0,
  max_score INT NOT NULL DEFAULT 0,
  percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  feedback TEXT,
  result_data JSONB NOT NULL DEFAULT '{}'::jsonb,
  graded_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_results_attempt
    FOREIGN KEY (attempt_id)
    REFERENCES attempts(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_results_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_results_test
    FOREIGN KEY (test_id)
    REFERENCES tests(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_results_attempt_id ON results(attempt_id);
CREATE INDEX idx_results_user_id ON results(user_id);
CREATE INDEX idx_results_test_id ON results(test_id);
CREATE INDEX idx_results_status ON results(status);
CREATE INDEX idx_results_result_data ON results USING GIN (result_data);

-- +goose Down
DROP TABLE IF EXISTS results;