package model

type Review struct {
	User         User   `json:"user,omitempty" structs:"user"`
	Comment      string `json:"comment,omitempty" structs:"comment"`
	Appreciation int    `json:"appreciation,omitempty" structs:"appreciation"`
}
