package ulasans

type RequestUlasan struct {
	DestinationId string  `json:"destination_id"`
	Message       string  `json:"message"`
	Rating        float32 `json:"rating"`
}

type ResponseUlasan struct {
	UserName        string `json:"user_name"`
	Message         string `json:"message"`
	Rating          float32 `json:"rating"`
}

type Ulasan struct {
	UserId        string  `json:"user_id"`
	DestinationId string  `json:"destination_id"`
	Message       string  `json:"message"`
	Rating        float32 `json:"rating"`
}
