package postgres

import (
	"context"
	"errors"
	"eventify-events/internal/models"
	"fmt"
	"reflect"
	"strings"
	"time"

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
var eventParticipantColumns = []string{"id", "user_id", "event_id", "is_owner", "can_edit_event", "can_manage_participants", "can_manage_checklist", "role", "status", "joined_at", "left_at"}
var eventChecklistColumns = []string{"id", "event_id", "title", "quantity", "unit", "is_purchased", "created_at"}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

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

func (r *EventRepository) UpdateEvent(ctx context.Context, params models.UpdateEventParams, id uuid.UUID) (models.Events, error) {
	updateData := make(map[string]any)

	v := reflect.ValueOf(params)
	t := reflect.TypeOf(params)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			tagName := fieldType.Tag.Get("json")
			columnName := strings.Split(tagName, ",")[0]

			if columnName != "" && columnName != "-" {
				updateData[columnName] = fieldValue.Elem().Interface()
			}
		}
	}

	if len(updateData) == 0 {
		return r.GetEvent(ctx, id)
	}

	sql, args, err := r.builder.
		Update("events").
		SetMap(updateData).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return models.Events{}, fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return models.Events{}, fmt.Errorf("failed to execute update: %w", err)
	}

	return r.GetEvent(ctx, id)
}

func (r *EventRepository) JoinEvent(ctx context.Context, userId uuid.UUID, code string) (bool, error) {
	sql, args, err := r.builder.Select("id").
		From("events").
		Where(squirrel.Eq{"event_code": code}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var eventId uuid.UUID
	err = r.db.QueryRow(ctx, sql, args...).Scan(&eventId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find event by code: %w", err)
	}

	_, joined, err := r.AddParticipant(ctx, userId, eventId)
	if err != nil {
		return false, fmt.Errorf("failed to add participant: %w", err)
	}
	return joined, nil
}

func (r *EventRepository) AddParticipant(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (uuid.UUID, bool, error) {
	p := models.EventParticipants{
		ID:       uuid.New(),
		UserID:   userId,
		EventID:  eventId,
		Status:   "confirmed",
		JoinedAt: time.Now(),
	}
	sql, args, err := r.builder.Insert("event_participants").
		Columns(eventParticipantColumns...).
		Values(p.Values()...).
		Suffix("ON CONFLICT (user_id, event_id) DO NOTHING").
		ToSql()

	if err != nil {
		return uuid.Nil, false, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("failed to execute insert: %w", err)
	}

	return eventId, res.RowsAffected() > 0, nil
}

func (r *EventRepository) RemoveParticipant(ctx context.Context, participantId uuid.UUID, eventId uuid.UUID) (bool, error) {
	sql, args, err := r.builder.
		Delete("event_participants").
		Where(squirrel.Eq{
			"user_id":  participantId,
			"event_id": eventId,
		}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("failed to execute delete: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

func (r *EventRepository) GetEventParticipants(ctx context.Context, eventId uuid.UUID) ([]models.EventParticipants, error) {
	sql, args, err := r.builder.Select(eventParticipantColumns...).
		From("event_participants").
		Where(squirrel.Eq{"event_id": eventId}).
		OrderBy("joined_at DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	participants, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.EventParticipants])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return participants, nil
}

func (r *EventRepository) CancelEvent(ctx context.Context, eventId uuid.UUID) (bool, error) {
	sql, args, err := r.builder.
		Update("events").
		Set("status", models.StatusCancelled).
		Where(squirrel.Eq{"id": eventId}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("failed to execute update: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

func (r *EventRepository) CreateInviteLink(ctx context.Context, eventId uuid.UUID, inviteType string, expiresAt *time.Time) (string, error) {
	var maxUses *int
	if inviteType == "single" {
		val := 1
		maxUses = &val
	}

	insertSQL, insertArgs, err := r.builder.Insert("event_invites").
		Columns("event_id", "invite_type", "expires_at", "max_uses").
		Values(eventId, inviteType, expiresAt, maxUses).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("failed to build insert: %w", err)
	}

	if _, err := r.db.Exec(ctx, insertSQL, insertArgs...); err != nil {
		return "", fmt.Errorf("failed to insert invite: %w", err)
	}

	selectSQL, selectArgs, err := r.builder.
		Select("event_code").
		From("events").
		Where(squirrel.Eq{"id": eventId}).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("failed to build select: %w", err)
	}

	var eventCode string
	err = r.db.QueryRow(ctx, selectSQL, selectArgs...).Scan(&eventCode)
	if err != nil {
		return "", fmt.Errorf("failed to get event code: %w", err)
	}

	return eventCode, nil
}

func (r *EventRepository) AddChecklistItem(ctx context.Context, e models.ChecklistItems) (uuid.UUID, error) {
	sql, args, err := r.builder.Insert("checklist_items").
		Columns(eventChecklistColumns...).
		Values(e.Values()...).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	var itemID uuid.UUID
	err = r.db.QueryRow(ctx, sql, args...).Scan(&itemID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert and scan id: %w", err)
	}
	return itemID, nil
}

func (r *EventRepository) GetEventChecklist(ctx context.Context, eventId uuid.UUID) ([]models.ChecklistItems, error) {
	sql, args, err := r.builder.Select(eventChecklistColumns...).
		From("checklist_items").
		Where(squirrel.Eq{"event_id": eventId}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get event_checklist: %w", err)
	}
	eventChecklist, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ChecklistItems])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}
	return eventChecklist, nil
}

func (r *EventRepository) RemoveChecklistItem(ctx context.Context, itemId uuid.UUID, eventId uuid.UUID) (bool, error) {
	sql, args, err := r.builder.
		Delete("checklist_items").
		Where(squirrel.Eq{
			"id":       itemId,
			"event_id": eventId,
		}).ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("failed to execute delete: %w", err)
	}
	return res.RowsAffected() > 0, nil
}

func (r *EventRepository) MarkItemPurchased(ctx context.Context, eventId uuid.UUID, itemId uuid.UUID, buyerId *uuid.UUID, isPurchased *bool) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if isPurchased != nil {
		sql, args, err := r.builder.
			Update("checklist_items").
			Set("is_purchased", *isPurchased).
			Where(squirrel.Eq{"id": itemId}).
			ToSql()

		if err != nil {
			return false, fmt.Errorf("failed to build update item query: %w", err)
		}

		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return false, fmt.Errorf("failed to update item status: %w", err)
		}
	}

	if buyerId != nil {
		sql, args, err := r.builder.
			Insert("checklist_assignments").
			Columns("checklist_item_id", "participant_id").
			Values(itemId, *buyerId).
			Suffix("ON CONFLICT (checklist_item_id, participant_id) DO NOTHING").
			ToSql()

		if err != nil {
			return false, fmt.Errorf("failed to build assignment query: %w", err)
		}

		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return false, fmt.Errorf("failed to assign buyer: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("failed to commit tx: %w", err)
	}

	return true, nil
}
