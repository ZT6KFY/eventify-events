package models

import (
	"time"

	"github.com/google/uuid"
)

type EventInvitesType string

const (
	TypeSingle    EventInvitesType = "single"
	TypeMulti     EventInvitesType = "multi"
	TypeUnlimited EventInvitesType = "unlimited"
)

type EventInvites struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	EventID    uuid.UUID        `json:"event_id" db:"event_id"`
	Token      string           `json:"token" db:"token"`
	InviteType EventInvitesType `json:"invite_type" db:"invite_type"`
	MaxUses    *int             `json:"max_uses" db:"max_uses"`
	UsedCount  int              `json:"used_count" db:"used_count"`
	ExpiresAt  *time.Time       `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at"`
}
