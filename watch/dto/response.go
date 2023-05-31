package dto

type CurrentUsersResponse struct {
	UsersCount int    `json:"users_count"`
	Users      []User `json:"users"`
}

type User struct {
	Name  string `json:"name" db:"name" binding:"required"`
	Email string `json:"email" db:"email" binding:"required"`
}

type CurrentVideoInfoResponse struct {
	Rooms []RoomInfo `json:"rooms"`
}

type RoomInfo struct {
	Name      string `json:"name"`
	Timing    int64  `json:"timing"`
	Url       string `json:"url"`
	UserCount int    `json:"users_count"`
	Duration  int    `json:"duration"`
	IsPrivate bool   `json:"is_private"`
	AdminName string `json:"admin_name"`
}
