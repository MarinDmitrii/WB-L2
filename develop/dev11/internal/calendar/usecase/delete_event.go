package usecase

import (
	"context"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type DeleteEventUseCase struct {
	eventRepository domain.Repository
}

func NewDeleteEventUseCase(
	eventRepository domain.Repository,
) *DeleteEventUseCase {
	return &DeleteEventUseCase{
		eventRepository: eventRepository,
	}
}

func (uc *DeleteEventUseCase) Execute(ctx context.Context, eventID int) error {
	return uc.eventRepository.DeleteEvent(ctx, eventID)
}
