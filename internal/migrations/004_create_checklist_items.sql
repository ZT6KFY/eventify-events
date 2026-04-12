-- +goose Up
CREATE TABLE IF NOT EXISTS checklist_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id     UUID         NOT NULL,
    title        VARCHAR(255) NOT NULL,
    quantity     INT          NOT NULL DEFAULT 1,
    unit         VARCHAR(50),
    is_purchased BOOLEAN      NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_checklist_event ON checklist_items (event_id);

-- +goose Down
DROP TABLE IF EXISTS checklist_items;
