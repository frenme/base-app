package entity

type CardEntity struct {
	ID            int
	Name          string
	Description   string
	FrontImageURL string
	BackImageURL  string
	NumUsersOwn   int
	NumUsersWish  int
	Status        string
}

func (CardEntity) TableName() string {
	return "cards"
}
