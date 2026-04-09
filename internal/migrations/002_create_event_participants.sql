-- +goose Up
CREATE TABLE IF NOT EXISTS event_participants (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID        NOT NULL,
    event_id                UUID        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    is_owner                BOOLEAN     NOT NULL DEFAULT false,
    can_edit_event          BOOLEAN     NOT NULL DEFAULT false,
    can_manage_participants BOOLEAN     NOT NULL DEFAULT false,
    can_manage_checklist    BOOLEAN     NOT NULL DEFAULT false,
    role                    VARCHAR(100),
    status                  VARCHAR(20) NOT NULL DEFAULT 'invited',
    joined_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    left_at                 TIMESTAMPTZ,

    CONSTRAINT check_participant_status CHECK (status IN ('invited', 'confirmed', 'declined', 'maybe')),
    UNIQUE (user_id, event_id)
);

CREATE INDEX idx_participants_event ON event_participants (event_id);
CREATE INDEX idx_participants_user  ON event_participants (user_id);

-- +goose Down
DROP TABLE IF EXISTS event_participants;
