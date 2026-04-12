package repository

import (
	"context"
	"eventify-events/internal/models"
	"time"

	"github.com/google/uuid"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event models.Events) error
	GetEvent(ctx context.Context, id uuid.UUID) (models.Events, error)
	ListUserEvents(ctx context.Context, userId uuid.UUID) ([]models.Events, error)
	ListEvents(ctx context.Context) ([]models.Events, error)
	UpdateEvent(ctx context.Context, params models.UpdateEventParams, id uuid.UUID) (models.Events, error)
	JoinEvent(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (uuid.UUID, bool, error)
	AddParticipant(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (uuid.UUID, bool, error)
	RemoveParticipant(ctx context.Context, participantId uuid.UUID, eventId uuid.UUID) (bool, error)
	GetEventParticipants(ctx context.Context, eventId uuid.UUID) ([]models.EventParticipants, error)
	CancelEvent(ctx context.Context, eventId uuid.UUID) (bool, error)
	CreateInviteLink(ctx context.Context, eventId uuid.UUID, inviteType string, expiresAt *time.Time) (string, error)
	AddChecklistItem(ctx context.Context, e models.ChecklistItems) (uuid.UUID, error)
	GetEventChecklist(ctx context.Context, eventId uuid.UUID) ([]models.ChecklistItems, error)
	RemoveChecklistItem(ctx context.Context, itemId uuid.UUID, eventId uuid.UUID) (bool, error)
	MarkItemPurchased(ctx context.Context, eventId uuid.UUID, itemId uuid.UUID, buyerId *uuid.UUID, isPurchased *bool) (bool, error)
}
