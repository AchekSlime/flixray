package dto

type Room struct {
	Name      string   `json:"name" binding:"required"`
	IsPrivate bool     `json:"is_private"`
	EmailList []string `json:"email_list"`
}

type AddUserRequest struct {
	RoomName string `json:"room_name" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
}

type DelRoomRequest struct {
	RoomName string `json:"room_name" binding:"required"`
}
