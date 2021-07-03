package model

type User struct {
	Name     string `json:"name" structs:"name"`
	Email    string `json:"email" structs:"email"`
	Password string `json:"password,omitempty" structs:"password"`
	Type     string `json:"type" structs:"type"`
}

type LoginData struct {
	Email    string `json:"email" structs:"email"`
	Password string `json:"password" structs:"password"`
}

type CheckUserData struct {
	Email string `json:"email" structs:"email"`
}
