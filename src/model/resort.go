package model

type Resort struct {
	Name        string   `json:"name,omitempty" structs:"name"`
	Description string   `json:"description,omitempty" structs:"description"`
	Price       float64  `json:"price,omitempty" structs:"price"`
	Specialties []string `json:"specialties,omitempty" structs:"specialties"`
	Amenities   []string `json:"amenities,omitempty" structs:"amenities"`
	Type        string   `json:"type,omitempty" structs:"type"`
	Location    string   `json:"location,omitempty" structs:"location"`
	Address     string   `json:"address,omitempty" structs:"address"`
	Reviews     []Review `json:"reviews,omitempty" structs:"reviews"`
	Images      []string `json:"images,omitempty" structs:"images"`
	Rating      float64  `json:"rating,omitempty" structs:"rating"`
	Tags        []string `json:"tags,omitempty" structs:"tags"`
}
