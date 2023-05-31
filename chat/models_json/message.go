package models_json

type Message struct {
	Token    string `json:"token"`
	Author   string `json:"author" binding:"required,gt=0,max=50"`
	RoomName string `json:"room_name" binding:"required,gt=0,max=50"`
	Message  string `json:"message" binding:"required,gt=0,max=200"`
}
