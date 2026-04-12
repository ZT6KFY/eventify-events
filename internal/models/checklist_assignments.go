package models

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistAssignments struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ChecklistItemId uuid.UUID `json:"checklist_item_id" db:"checklist_item_id"`
	ParticipantId   uuid.UUID `json:"participant_id" db:"participant_id"`
	Quantity        int       `json:"quantity" db:"quantity"`
	AssignedAt      time.Time `json:"assigned_at" db:"assigned_at"`
}
