package globe

import (
	"fmt"

	"github.com/google/uuid"
)

type ID string

func NewID() ID {
	return ID(uuid.New().String())
}

type PID struct {
	Address string
	ID      ID
}

type pidoption func(*PID)

func NewPID(options ...pidoption) PID {
	id := NewID()
	pid := PID{
		ID:      id,
		Address: fmt.Sprintf("actor://%s", id),
	}
	for _, opt := range options {
		opt(&pid)
	}
	return pid
}
