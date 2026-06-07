-- +goose Up
CREATE TABLE IF NOT EXISTS sections (
  id UUID PRIMARY KEY,
  test_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  section_code VARCHAR(50) NOT NULL,
  display_order INT NOT NULL,
  duration_minutes INT NOT NULL,
  total_questions INT NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_sections_test
    FOREIGN KEY (test_id)
    REFERENCES tests(id)
    ON DELETE CASCADE,

  CONSTRAINT uq_sections_test_code
    UNIQUE (test_id, section_code)
);

CREATE INDEX idx_sections_test_id
ON sections(test_id);

CREATE INDEX idx_sections_display_order
ON sections(display_order);

CREATE INDEX idx_sections_is_active
ON sections(is_active);

-- +goose Down
DROP TABLE IF EXISTS sections;