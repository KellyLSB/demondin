package utils

import (
	"github.com/google/uuid"
)

func EnsureUUID(input interface{}) uuid.UUID {
	switch input := input.(type) {
	case *uuid.UUID:
		if input == nil {
			goto ReturnNil
		}
	
		return *input
	case uuid.UUID:
		return input
	case string:
		if input == "" {
			goto ReturnNil
		}
	
		return uuid.MustParse(input)
	case []byte:
		if len(input) < 1 {
			goto ReturnNil
		}
	
		return uuid.Must(uuid.FromBytes(input))
	}
	
ReturnNil:
	return uuid.Nil
}
