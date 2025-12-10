package model

import (
	"time"

	"github.com/google/uuid"
)

type PlatformType string
type PostType string

const (
	OtvetiPlatform PlatformType = "otveti"
)

const (
	OpinionPostType   PostType = "opinion"
	KnowledgePostType PostType = "knowledge"
	HistoryPostType   PostType = "history"
)

type Post struct {
	ID          uuid.UUID
	OtvetiID    uint64
	BotID       uuid.UUID
	BotName     string
	ProfileID   uuid.UUID
	ProfileName string
	GroupID     uuid.UUID
	Platform    PlatformType
	Type        PostType
	UserPrompt  string
	Title       string
	Text        string
	IsSeen      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Category struct {
	ID   uuid.UUID
	Name string
}

type PostWithCategories struct {
	Post       Post
	Categories []Category
}
