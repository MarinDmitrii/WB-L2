package adapters

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MarinDmitrii/WB-L2/develop/dev11/internal/calendar/domain"
)

type CacheEventRepository struct {
	cache         map[int]domain.Event
	autoIncrement int
	maxSize       int
	mu            *sync.RWMutex
}

func NewCacheEventRepository(maxSize int) *CacheEventRepository {
	return &CacheEventRepository{
		cache:         make(map[int]domain.Event, maxSize),
		autoIncrement: 1,
		maxSize:       maxSize,
		mu:            &sync.RWMutex{},
	}
}

func (r *CacheEventRepository) CreateEvent(ctx context.Context, domainEvent domain.Event) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.autoIncrement == r.maxSize {
		r.autoIncrement = 1
	}

	if len(r.cache) == r.maxSize {
		delete(r.cache, r.autoIncrement)
	}

	if _, ok := r.cache[domainEvent.ID]; !ok {
		domainEvent.ID = r.autoIncrement
		r.autoIncrement++
	}

	r.cache[domainEvent.ID] = domainEvent

	return domainEvent.ID, nil
}

func (r *CacheEventRepository) GetEventByID(ctx context.Context, eventID int) (domain.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if event, ok := r.cache[eventID]; !ok {
		return domain.Event{}, fmt.Errorf("Error: can't find event\n")
	} else {
		return event, nil
	}
}

func (r *CacheEventRepository) UpdateEvent(ctx context.Context, updatedEvent domain.Event) error {
	_, err := r.GetEventByID(ctx, updatedEvent.ID)
	if err != nil {
		return err
	}

	_, err = r.CreateEvent(ctx, updatedEvent)
	if err != nil {
		return err
	}

	return nil
}

func (r *CacheEventRepository) DeleteEvent(ctx context.Context, eventID int) error {
	_, err := r.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	r.mu.Lock()
	delete(r.cache, eventID)
	r.mu.Unlock()

	return nil
}

func (r *CacheEventRepository) GetEventsForDay(ctx context.Context, date time.Time) ([]domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	eventsForDay := make([]domain.Event, 0, 5)

	for _, v := range r.cache {
		if v.Date.Year() == date.Year() && v.Date.Month() == date.Month() && v.Date.Day() == date.Day() {
			eventsForDay = append(eventsForDay, v)
		}
	}

	return eventsForDay, nil
}

func (r *CacheEventRepository) GetEventsForWeek(ctx context.Context, date time.Time) ([]domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	eventsForWeek := make([]domain.Event, 0, 10)

	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	// Определяем пограничные дни недели для корректной работы функций time.After() и time.Before()
	startOfWeek := date.AddDate(0, 0, -weekday)
	endOfWeek := startOfWeek.AddDate(0, 0, 8)

	for _, v := range r.cache {
		if v.Date.After(startOfWeek) && v.Date.Before(endOfWeek) {
			eventsForWeek = append(eventsForWeek, v)
		}
	}

	return eventsForWeek, nil
}

func (r *CacheEventRepository) GetEventsForMonth(ctx context.Context, date time.Time) ([]domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	eventsForMonth := make([]domain.Event, 0, 20)

	for _, v := range r.cache {
		if v.Date.Year() == date.Year() && v.Date.Month() == date.Month() {
			eventsForMonth = append(eventsForMonth, v)
		}
	}

	return eventsForMonth, nil
}
