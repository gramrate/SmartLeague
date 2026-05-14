package series

import "github.com/google/uuid"

func maybeRequesterID(v any) *uuid.UUID {
	id, ok := v.(uuid.UUID)
	if !ok || id == uuid.Nil {
		return nil
	}
	return &id
}
