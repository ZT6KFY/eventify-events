package postgres_test

import (
	"context"
	"testing"
	"time"

	"eventify-events/internal/models"
	"eventify-events/internal/repository/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Ptr[T any](v T) *T { return &v }

func TestEventRepository_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	dbURL := "postgres://user:password@localhost:5432/postgres?sslmode=disable"

	pool, err := postgres.NewPool(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgres.NewEventRepository(pool)

	_, err = pool.Exec(ctx, "DELETE FROM events")
	require.NoError(t, err)

	eventID := uuid.New()
	testEvent := models.Events{
		ID:        eventID,
		CreatorID: uuid.New(),
		Title:     "Test Event",
		StartsAt:  time.Now().Add(time.Hour).Truncate(time.Second),
		Status:    models.StatusDraft,
	}

	t.Run("Create event", func(t *testing.T) {
		err := repo.CreateEvent(ctx, testEvent)
		assert.NoError(t, err)
	})

	t.Run("Get existing event", func(t *testing.T) {
		found, err := repo.GetEvent(ctx, eventID)
		assert.NoError(t, err)
		assert.Equal(t, testEvent.ID, found.ID)
		assert.Equal(t, testEvent.Title, found.Title)
	})

	t.Run("Get non-existent event", func(t *testing.T) {
		_, err := repo.GetEvent(ctx, uuid.New())
		assert.Error(t, err)
	})

	t.Run("Get all events", func(t *testing.T) {
		_, err := repo.ListEvents(ctx)
		assert.NoError(t, err)
	})
}
