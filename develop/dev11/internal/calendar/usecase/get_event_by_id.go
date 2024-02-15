package usecase

import (
	"context"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type GetEventByIDUseCase struct {
	eventRepository domain.Repository
}

func NewGetEventByIDUseCase(
	eventRepository domain.Repository,
) *GetEventByIDUseCase {
	return &GetEventByIDUseCase{
		eventRepository: eventRepository,
	}
}

func (uc *GetEventByIDUseCase) Execute(ctx context.Context, eventID int) (domain.Event, error) {
	return uc.eventRepository.GetEventByID(ctx, eventID)
}
