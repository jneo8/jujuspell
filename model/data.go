package model

import (
	"github.com/google/uuid"
)

// Send to bubbletea model to refresh data
type RefreshMsg interface {
	GetJobID() uuid.UUID
}

type ResourceType string
type Filter string

type DataRefresher interface {
	Fetch(filter Filter) error
	GetRefreshMsg() RefreshMsg
}

type QueryJob struct {
	ID           uuid.UUID
	ResourceType ResourceType
	Filter       Filter
}
