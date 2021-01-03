package utils

import(
	"github.com/google/uuid"
)

func GetArrayOfUUIDFromMapOfUUID(
	IDs map[uuid.UUID]bool,
)[]uuid.UUID {
	res := make([]uuid.UUID, 0, len(IDs))
	for id := range IDs {
		res = append(res, id)
	}
	return res
}

