package favorites

import "time"

type Favorite struct {
	ID            string
	UserId        string
	DestinationId string
	CreatedAt     time.Time
}

func (Favorite) TableName() string {
	return "tbl_favorite"
}
