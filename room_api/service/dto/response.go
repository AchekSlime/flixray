package dto

type AvailableRoom struct {
	Name      string         `json:"name"`
	IsPrivate bool           `json:"is_private"`
	Admin     *AvailableUser `json:"admin"`
}

type AvailableUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
