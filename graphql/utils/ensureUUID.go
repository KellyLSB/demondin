package utils

import (
	"github.com/google/uuid"
)

func EnsureUUID(input interface{}) uuid.UUID {
	switch input := input.(type) {
	case *uuid.UUID:
		if input == nil {
			return uuid.Nil
		}
	
		return *input
	case uuid.UUID:
		return input
	case string:
		return uuid.MustParse(input)
	case []byte:
		return uuid.Must(uuid.FromBytes(input))
	}
	
	return uuid.Nil
}
