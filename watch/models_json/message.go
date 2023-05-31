package models_json

type Event struct {
	Author string `json:"author" binding:"required,gt=0,max=50"`
	Type   string `json:"type" binding:"required,gt=0,max=200"`
	Timing int64  `json:"timing" binding:"required,gt=0,max=200"`
	Status string `json:"status"`
}
