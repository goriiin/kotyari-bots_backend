package model

import (
	"time"

	"github.com/google/uuid"
)

type PlatformType string

const (
	OtvetiPlatform PlatformType = "otveti"
)

type PostType string

const (
	OpinionPostType   PostType = "opinion"
	KnowledgePostType PostType = "knowledge"
	HistoryPostType   PostType = "history"
)

type Post struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	OtvetiID    uint64
	BotID       uuid.UUID
	BotName     string
	ProfileID   uuid.UUID
	ProfileName string
	GroupID     uuid.UUID
	UserPrompt  string
	Platform    PlatformType
	Type        PostType
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
