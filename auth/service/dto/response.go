package dto

type RegUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegMailResponse struct {
	Code int64 `json:"code"`
}
