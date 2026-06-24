package util

import "github.com/google/uuid"

type NewUUID func() uuid.UUID

func MakeNewUUID() uuid.UUID {

	return uuid.New()
}
