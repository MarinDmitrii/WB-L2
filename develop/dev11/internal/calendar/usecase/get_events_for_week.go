package usecase

import (
	"context"
	"time"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type GetEventsForWeekUseCase struct {
	eventRepository domain.Repository
}

func NewGetEventsForWeekUseCase(
	eventRepository domain.Repository,
) *GetEventsForWeekUseCase {
	return &GetEventsForWeekUseCase{
		eventRepository: eventRepository,
	}
}

func (uc *GetEventsForWeekUseCase) Execute(ctx context.Context, date time.Time) ([]domain.Event, error) {
	return uc.eventRepository.GetEventsForWeek(ctx, date)
}
