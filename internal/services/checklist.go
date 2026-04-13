package services

import (
	"context"
	"eventify-events/internal/models"
	"eventify-events/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChecklistService struct {
	repo repository.EventRepository
}

func NewChecklistService(repo repository.EventRepository) *ChecklistService {
	return &ChecklistService{repo: repo}
}

func (s *ChecklistService) AddChecklistItem(ctx context.Context, callerID uuid.UUID, eventID uuid.UUID, title string, quantity int, unit string) (uuid.UUID, error)

func (s *ChecklistService) RemoveChecklistItem(ctx context.Context, callerID uuid.UUID, eventID uuid.UUID, itemID uuid.UUID) (bool, error)

func (s *ChecklistService) MarkItemPurchased(ctx context.Context, callerID uuid.UUID, eventID uuid.UUID, itemID uuid.UUID, buyerID *uuid.UUID, isPurchased *bool) (bool, error)

func (s *ChecklistService) GetEventChecklist(ctx context.Context, callerID uuid.UUID, eventID uuid.UUID) ([]models.ChecklistItems, error)
