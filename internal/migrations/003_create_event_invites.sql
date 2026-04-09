-- +goose Up
CREATE TABLE IF NOT EXISTS event_invites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    token       VARCHAR(64) NOT NULL UNIQUE DEFAULT substr(md5(random()::text || now()::text), 1, 16),
    invite_type VARCHAR(20) NOT NULL DEFAULT 'unlimited',
    max_uses    INT,
    used_count  INT         NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT check_invite_type CHECK (invite_type IN ('single', 'multi', 'unlimited'))
);

CREATE INDEX idx_invites_event ON event_invites (event_id);
CREATE INDEX idx_invites_token ON event_invites (token);

-- +goose Down
DROP TABLE IF EXISTS event_invites;
