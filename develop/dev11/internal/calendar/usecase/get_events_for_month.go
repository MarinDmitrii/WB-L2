package usecase

import (
	"context"
	"time"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type GetEventsForMonthUseCase struct {
	eventRepository domain.Repository
}

func NewGetEventsForMonthUseCase(
	eventRepository domain.Repository,
) *GetEventsForMonthUseCase {
	return &GetEventsForMonthUseCase{
		eventRepository: eventRepository,
	}
}

func (uc *GetEventsForMonthUseCase) Execute(ctx context.Context, date time.Time) ([]domain.Event, error) {
	return uc.eventRepository.GetEventsForMonth(ctx, date)
}
