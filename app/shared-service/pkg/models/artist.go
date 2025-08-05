package models

import "time"

type Artist struct {
	ID           int
	AgencyID     int
	Name         string
	DebutDate    time.Time
	NumFollowers int
	ImageUrl     *string
	Status       ArtistStatus
}

type ArtistStatus string

const (
	Active    ArtistStatus = "active"
	Suspended ArtistStatus = "suspended"
	Disband   ArtistStatus = "disband"
)
