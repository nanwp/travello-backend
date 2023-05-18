package destinations

type Destination struct {
	ID          string  `json:"_id"`
	Nama        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      float32 `json:"rating"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type DestinationResponse struct {
	ID          string  `json:"id,omitempty"`
	Nama        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Location    string  `json:"location,omitempty"`
	Category    string  `json:"category,omitempty"`
	Image       string  `json:"image,omitempty"`
	Rating      float32 `json:"rating,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}
