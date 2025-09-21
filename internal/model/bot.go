package model

import (
	"time"

	"github.com/google/uuid"
)

type Bot struct {
	ID                 uuid.UUID
	Name               string
	Email              string
	SystemPrompt       string
	ModerationRequired bool
	ProfileIDs         []uuid.UUID
	ProfilesCount      int
	CreatedAt          time.Time
	UpdateAt           time.Time
}
