package entity

type AgencyEntity struct {
	ID       int
	Name     string
	ImageUrl *string
	Status   string
}

func (AgencyEntity) TableName() string {
	return "agencies"
}
