package model

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	ID          uuid.UUID
	Name        string
	Email       string
	SystemPromt string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
