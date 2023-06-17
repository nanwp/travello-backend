package destinations

import "github.com/nanwp/travello/models/ulasans"

type Destination struct {
	ID          string                   `json:"_id"`
	Nama        string                   `json:"name" binding:"required"`
	Description string                   `json:"description" binding:"required"`
	Location    string                   `json:"location" binding:"required"`
	Category    string                   `json:"category"`
	Image       []string                 `json:"image"`
	Rating      float32                  `json:"rating"`
	CountUlasan int                      `json:"jumlah_ulasan"`
	Ulasan      []ulasans.ResponseUlasan `json:"ulasan"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

type DestinationCreate struct {
	Nama        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Location    string   `json:"location" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Image       []string `json:"image"`
}

type DestinationUpdate struct {
	Nama        string   `json:"name"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Category    string   `json:"category"`
	Image       []string `json:"image"`
}

type DestinationResponse struct {
	ID          string                   `json:"id,omitempty"`
	Nama        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Location    string                   `json:"location,omitempty"`
	Category    string                   `json:"category,omitempty"`
	Image       []string                 `json:"image,omitempty"`
	Rating      float32                  `json:"rating"`
	CountUlasan int                      `json:"jumlah_ulasan"`
	Ulasan      []ulasans.ResponseUlasan `json:"ulasan"`
	CreatedAt   string                   `json:"created_at,omitempty"`
	UpdatedAt   string                   `json:"updated_at,omitempty"`
}
