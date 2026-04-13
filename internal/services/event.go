package services

import (
	"context"
	"eventify-events/internal/models"
	"eventify-events/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) checkPermission(ctx context.Context, userID, eventID uuid.UUID, permission string) error {
	participants, err := s.repo.GetEventParticipants(ctx, eventID) // Эффективнее сделать GetEventParticapant в репозитории, чтобы не делать лишних запросов
	if err != nil {
		return fmt.Errorf("checkPermission: %w", err)
	}

	for _, p := range participants {
		if p.UserID != userID {
			continue
		}
		if p.IsOwner {
			return nil
		}
		switch permission {
		case "can_edit_event":
			if p.CanEditEvent {
				return nil
			}
		case "can_manage_participants":
			if p.CanManageParticipants {
				return nil
			}
		case "can_manage_checklist":
			if p.CanManageChecklist {
				return nil
			}
		}
		return fmt.Errorf("permission denied: user %s lacks %s", userID, permission)
	}
	return fmt.Errorf("permission denied: user %s is not a participant", userID)
}

func (s *EventService) GetEvent(ctx context.Context, uuid uuid.UUID) (*models.Events, error) {
	event, err := s.repo.GetEvent(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("Service.GetEvent : %w", err)
	}
	return &event, nil
}

func (s *EventService) ListEvents(ctx context.Context) ([]models.Events, error) { // Возвращает все неотмененные и публичные ивенты
	events, err := s.repo.ListEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Service.ListEvents : %w", err)
	}
	result := make([]models.Events, 0, len(events))
	for _, e := range events {
		if e.Status != models.StatusCancelled && !e.IsPrivate {
			result = append(result, e)
		}
	}
	return result, nil
}


func (s *EventService) ListUserEvents(ctx context.Context, userId uuid.UUID) ([]models.Events, error) { // Пока возвращает только ивенты, где юзер - создатель (т.е не вовзращет, где он участник)
	events, err := s.repo.ListUserEvents(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("Service.ListUserEvents : %w", err)
	}
	return events, nil
}

func (s *EventService) JoinEvent(ctx context.Context, userId uuid.UUID, code string) (bool, error) {
	event, err := s.repo.GetEventByCode(ctx, code)
	if err != nil {
		return false, fmt.Errorf("Service.GetEventByCode : %w", err)
	}
	particapants, err := s.repo.GetEventParticipants(ctx, event.ID)
	if err != nil {
		return false, fmt.Errorf("Service.GetEventParticipants : %w", err)
	}


	if event.MaxParticipants != nil && *event.MaxParticipants != 0 && len(particapants) >= *event.MaxParticipants {
		return false, fmt.Errorf("event with code %s is full", code)
	}

	// Проверка на статус (не cancelled и не completed)
	if event.Status == models.StatusCancelled || event.Status == models.StatusCompleted {
		return false, fmt.Errorf("event with code %s is %s", code, event.Status)
	}
	_, joined, err := s.repo.JoinEvent(ctx, userId, event.ID)
	if err != nil {
		return false, fmt.Errorf("Service.JoinEvent : %w", err)
	}
	return joined, nil
}

func (s *EventService) RemoveParticipant(ctx context.Context, callerId uuid.UUID, participantId uuid.UUID, eventId uuid.UUID) (bool, error) {
	err := s.checkPermission(ctx, callerId, eventId, "can_manage_participants")
	if err != nil {
		return false, fmt.Errorf("Service.RemoveParticipant : %w", err)
	}


	particapants, err := s.repo.GetEventParticipants(ctx, eventId)
	if err != nil {
		return false, fmt.Errorf("Service.RemoveParticipant : %w", err)
	}

	var particapantFound bool
	for _, p := range particapants {
		if p.UserID == participantId {
			// Проверка на удаление создателя
			if p.IsOwner {
				return false, fmt.Errorf("Service.RemoveParticipant : can't remove creator")
			}
			particapantFound = true
			break
		}
	}

	// Проверка на наличие участника
	if !particapantFound {
		return false, fmt.Errorf("Service.RemoveParticipant : participant not found")
	}

	

	removed, err := s.repo.RemoveParticipant(ctx, participantId, eventId)
	if err != nil {
		return false, fmt.Errorf("Service.RemoveParticipant : %w", err)
	}
	return removed, nil
}