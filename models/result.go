package models

import (
	"time"

	"gorm.io/datatypes"
)

type Status string

const (
	StatusQueued     Status = "Queued"
	StatusInProgress Status = "In Progress"
	StatusSuccess    Status = "Success"
	StatusFailure    Status = "Failure"
)

type Result struct {
	ID         string         `json:"scanId" gorm:"default:uuid_generate_v3()"`
	Status     Status         `json:"status"`
	RepoID     string         `json:"-"`
	RepoName   string         `json:"repositoryName" gorm:"->"`
	RepoUrl    string         `json:"repositoryUrl" gorm:"->"`
	Findings   datatypes.JSON `json:"findings,omitempty" gorm:"type:jsonb"`
	QueuedAt   time.Time      `json:"queuedAt"`
	ScanningAt time.Time      `json:"scanningAt,omitempty"`
	FinishedAt time.Time      `json:"finishedAt,omitempty"`
}
