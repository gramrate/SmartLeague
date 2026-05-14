package series

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func maybeRequesterID(v any) *uuid.UUID {
	id, ok := v.(uuid.UUID)
	if !ok || id == uuid.Nil {
		return nil
	}
	return &id
}

func parseDateTimeLocal(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty datetime")
	}
	// Front sends datetime-local as "2006-01-02T15:04" without timezone.
	return time.Parse("2006-01-02T15:04", raw)
}

func parseDateTimeInput(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty datetime")
	}
	if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
		return parsed, nil
	}
	return parseDateTimeLocal(raw)
}
