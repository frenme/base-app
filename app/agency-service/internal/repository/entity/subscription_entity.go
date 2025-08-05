package entity

type SubscriptionEntity struct {
	ID       int
	UserID   int
	ArtistID int
}

func (SubscriptionEntity) TableName() string {
	return "subscriptions"
}
