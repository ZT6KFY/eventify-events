package services

import (
	"context"
	"eventify-events/internal/models"
	"eventify-events/internal/repository"
	"fmt"
	"github.com/google/uuid"
)
func checkPermission(ctx context.Context, repo repository.EventRepository, userID, eventID uuid.UUID, permission string) error {
	participant, err := repo.GetParticipant(ctx, userID,eventID) 
	if err != nil {
		return fmt.Errorf("checkPermission: %w", err)
	}

	if participant.IsOwner {
		return nil
	}
	switch permission {
	case "can_edit_event":
		if participant.CanEditEvent {
			return nil
		}
	case "can_manage_participants":
		if participant.CanManageParticipants {
			return nil
		}
	case "can_manage_checklist":
		if participant.CanManageChecklist {
			return nil
		}
	}
	return fmt.Errorf("permission denied: user %s lacks %s", userID, permission)
}

var availableStatuses = map[models.EventStatus]bool{
    models.StatusDraft:  true,
    models.StatusActive: true,
}
func checkEventStatus(status models.EventStatus) error {
	if !availableStatuses[status] {
		return fmt.Errorf("event status %s is not available", status)
	}
	return nil
}