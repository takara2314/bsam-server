package domain

import "time"

type Race struct {
	AssociationID string      `json:"association_id"`
	Name          string      `json:"name"`
	Started       bool        `json:"started"`
	StartedAt     time.Time   `json:"started_at"`
	FinishedAt    time.Time   `json:"finished_at"`
	Association   Association `json:"association"`
	AthleteIDs    []string    `json:"athlete_ids"`
	MarkIDs       []string    `json:"mark_ids"`
}
