package model

import (
	"github.com/google/uuid"
	"time"
)

type Bot struct {
	ID                 uuid.UUID
	Name               string
	Email              string
	SystemPrompt       string
	ModerationRequired bool
	ProfileIDs         []uuid.UUID
	AutoPublish        bool
	CreatedAt          time.Time
	UpdateAt           time.Time
}
