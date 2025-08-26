package models

type Card struct {
	ID            int
	Name          string
	Description   string
	FrontImageURL string
	BackImageURL  string
	NumUsersOwn   int
	NumUsersWish  int
	Status        CollectStatus
}

type CardSet struct {
	ID    int
	Name  string
	Cards []Card
}
