package usecase

import (
	"context"
	"time"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type GetEventsForDayUseCase struct {
	eventRepository domain.Repository
}

func NewGetEventsForDayUseCase(
	eventRepository domain.Repository,
) *GetEventsForDayUseCase {
	return &GetEventsForDayUseCase{
		eventRepository: eventRepository,
	}
}

func (uc *GetEventsForDayUseCase) Execute(ctx context.Context, date time.Time) ([]domain.Event, error) {
	return uc.eventRepository.GetEventsForDay(ctx, date)
}
