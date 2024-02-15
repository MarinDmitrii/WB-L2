package domain

import "time"

type Event struct {
	ID          int
	UserID      int
	Date        time.Time
	Description string
}
