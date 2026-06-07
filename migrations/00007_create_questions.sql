-- +goose Up
CREATE TABLE IF NOT EXISTS questions (
  id UUID PRIMARY KEY,
  section_id UUID NOT NULL,
  question_type VARCHAR(50) NOT NULL,
  question_text TEXT NOT NULL,
  question_data JSONB NOT NULL DEFAULT '{}'::jsonb,
  correct_answer TEXT,
  explanation TEXT,
  points INT NOT NULL DEFAULT 1,
  display_order INT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_questions_section
    FOREIGN KEY (section_id)
    REFERENCES sections(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_questions_section_id
ON questions(section_id);

CREATE INDEX idx_questions_question_type
ON questions(question_type);

CREATE INDEX idx_questions_display_order
ON questions(display_order);

CREATE INDEX idx_questions_is_active
ON questions(is_active);

CREATE INDEX idx_questions_question_data
ON questions USING GIN (question_data);

-- +goose Down
DROP TABLE IF EXISTS questions;