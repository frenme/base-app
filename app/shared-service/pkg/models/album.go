package models

type Album struct {
	ID            int
	Name          string
	ImageUrl      string
	ReleaseDate   int
	AlbumVersions []AlbumVersion
	CardSets      []CardSet
}

type AlbumVersion struct {
	ID       int
	Name     string
	ImageUrl string
	Status   CollectStatus
}
