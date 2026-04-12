-- +goose Up
CREATE TABLE IF NOT EXISTS checklist_assignments (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    checklist_item_id UUID NOT NULL REFERENCES checklist_items(id) ON DELETE CASCADE,
    participant_id    UUID NOT NULL,
    quantity          INT  NOT NULL DEFAULT 1,
    assigned_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (checklist_item_id, participant_id)
);

CREATE INDEX idx_checklist_assign_item ON checklist_assignments (checklist_item_id);
CREATE INDEX idx_checklist_assign_part ON checklist_assignments (participant_id);

-- +goose Down
DROP TABLE IF EXISTS checklist_assignments;
