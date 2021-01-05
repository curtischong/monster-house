package utils

import (
	"github.com/google/uuid"
)

// GetArrayOfUUIDFromMapOfUUID returns an array of UUID
// for every UUID in the given map
func GetArrayOfUUIDFromMapOfUUID(
	IDs map[uuid.UUID]bool,
) []uuid.UUID {
	res := make([]uuid.UUID, 0, len(IDs))
	for id := range IDs {
		res = append(res, id)
	}
	return res
}
