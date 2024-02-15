package usecase

import (
	"context"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type CreateEventUseCase struct {
	eventRepository domain.Repository
}

func NewCreateEventUseCase(
	eventRepository domain.Repository,
) *CreateEventUseCase {
	return &CreateEventUseCase{
		eventRepository: eventRepository,
	}
}

func (uc *CreateEventUseCase) Execute(ctx context.Context, event domain.Event) (int, error) {
	return uc.eventRepository.CreateEvent(ctx, event)
}
