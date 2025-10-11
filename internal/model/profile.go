package model

import "github.com/google/uuid"

type Profile struct {
	ID          uuid.UUID
	Name        string
	Email       string
	SystemPromt string
}
