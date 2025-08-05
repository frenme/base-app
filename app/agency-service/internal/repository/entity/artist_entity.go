package entity

import "time"

type ArtistEntity struct {
	ID        int
	AgencyID  int
	Name      string
	DebutDate time.Time
	ImageUrl  *string
	Status    string
}

func (ArtistEntity) TableName() string {
	return "artists"
}
