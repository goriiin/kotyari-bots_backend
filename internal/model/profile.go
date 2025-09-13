package model

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID          uuid.UUID
	Name        string
	Email       string
	SystemPromt string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
