package ulasans

import "github.com/nanwp/travello/models/users"

type RequestUlasan struct {
	DestinationId string  `json:"destination_id"`
	Message       string  `json:"message"`
	Rating        float32 `json:"rating"`
}

type ResponseUlasan struct {
	UserName string  `json:"user_name"`
	Message  string  `json:"message"`
	Rating   float32 `json:"rating"`
}

type Ulasan struct {
	ID            string
	UserId        string
	User          users.User `gorm:"foreigenKey:UserId"`
	DestinationId string
	Message       string
	Rating        float32
}

func (Ulasan) TableName() string {
	return "tbl_ulasan"
}
