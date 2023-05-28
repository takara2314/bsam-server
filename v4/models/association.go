package models

import "time"

type AssociationsGETAllRes struct {
	Assocs []Association `json:"associations"`
}

type Association struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	TokenIAT time.Time `json:"-"`
	TokenEXP time.Time `json:"-"`
	Lat      float64   `json:"latitude"`
	Lng      float64   `json:"longitude"`
	RaceName string    `json:"race_name"`
}
