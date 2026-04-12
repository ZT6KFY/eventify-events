package models

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistItems struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	Title       string    `json:"title" db:"title"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Unit        *string   `json:"unit,omitempty" db:"unit"`
	IsPurchased bool      `json:"is_purchased" db:"is_purchased"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func (c ChecklistItems) Values() []any {
	return []any{
		c.ID, c.EventID, c.Title, c.Quantity, c.Unit,
		c.IsPurchased, c.CreatedAt,
	}
}
