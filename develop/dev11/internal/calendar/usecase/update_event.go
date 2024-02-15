package usecase

import (
	"context"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type UpdateEventUseCase struct {
	eventRepository domain.Repository
}

func NewUpdateEventUseCase(
	eventRepository domain.Repository,
) *UpdateEventUseCase {
	return &UpdateEventUseCase{
		eventRepository: eventRepository,
	}
}

func (uc *UpdateEventUseCase) Execute(ctx context.Context, updatedEvent domain.Event) error {
	return uc.eventRepository.UpdateEvent(ctx, updatedEvent)
}
