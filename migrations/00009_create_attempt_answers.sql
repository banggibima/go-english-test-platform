-- +goose Up
CREATE TABLE IF NOT EXISTS attempt_answers (
  id UUID PRIMARY KEY,
  attempt_id UUID NOT NULL,
  question_id UUID NOT NULL,
  answer_text TEXT,
  answer_data JSONB NOT NULL DEFAULT '{}'::jsonb,
  is_correct BOOLEAN,
  points_awarded INT NOT NULL DEFAULT 0,
  answered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_attempt_answers_attempt
    FOREIGN KEY (attempt_id)
    REFERENCES attempts(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_attempt_answers_question
    FOREIGN KEY (question_id)
    REFERENCES questions(id)
    ON DELETE CASCADE,

  CONSTRAINT uq_attempt_answers_attempt_question
    UNIQUE (attempt_id, question_id)
);

CREATE INDEX idx_attempt_answers_attempt_id ON attempt_answers(attempt_id);
CREATE INDEX idx_attempt_answers_question_id ON attempt_answers(question_id);
CREATE INDEX idx_attempt_answers_answer_data ON attempt_answers USING GIN (answer_data);

-- +goose Down
DROP TABLE IF EXISTS attempt_answers;