package postgres

import (
	"context"
	"eventify-events/internal/models"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

var eventColumns = []string{"id", "creator_id", "is_private", "title", "description", "starts_at", "duration", "location_name", "location_coords", "max_participants", "status", "event_code", "created_at", "updated_at"}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// create
func (r *EventRepository) CreateEvent(ctx context.Context, e models.Events) error {
	sql, args, err := r.builder.Insert("events").
		Columns(eventColumns...).
		Values(e.Values()...).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query %w", err)
	}
	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

// read
func (r *EventRepository) GetEvent(ctx context.Context, id uuid.UUID) (models.Events, error) {
	sql, args, err := r.builder.Select(eventColumns...).
		From("events").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return models.Events{}, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return models.Events{}, fmt.Errorf("failed to get event: %w", err)
	}

	event, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Events])
	if err != nil {
		return models.Events{}, fmt.Errorf("failed to collect row: %w", err)
	}
	return event, nil
}

func (r *EventRepository) ListUserEvents(ctx context.Context, userId uuid.UUID) ([]models.Events, error) {
	sql, args, err := r.builder.Select(eventColumns...).
		From("events").
		Where(squirrel.Eq{"creator_id": userId}).
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Events])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row")
	}
	return events, nil
}

func (r *EventRepository) ListEvents(ctx context.Context) ([]models.Events, error) {
	sql, args, err := r.builder.Select(eventColumns...).
		From("events").
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Events])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}
	return events, nil
}
