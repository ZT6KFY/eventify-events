-- +goose Up
CREATE TYPE event_status AS ENUM ('draft', 'active', 'cancelled', 'completed');

CREATE TABLE IF NOT EXISTS events (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id       UUID         NOT NULL,
    is_private       BOOLEAN      NOT NULL DEFAULT false,
    title            VARCHAR(255) NOT NULL,
    description      TEXT,
    starts_at        TIMESTAMPTZ  NOT NULL,
    duration         INTERVAL,
    location_name    VARCHAR(255),
    location_coords  POINT,
    max_participants INT,
    status           event_status NOT NULL DEFAULT 'draft',
    event_code       VARCHAR(32)  NOT NULL UNIQUE DEFAULT substr(md5(random()::text), 1, 10),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_creator   ON events (creator_id);
CREATE INDEX idx_events_starts_at ON events (starts_at);
CREATE INDEX idx_events_status    ON events (status);
CREATE INDEX idx_events_code      ON events (event_code);
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER set_updated_at_events
    BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at_events ON events;
DROP FUNCTION IF EXISTS trigger_set_updated_at;
DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS event_status;