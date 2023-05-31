package dto

type ConfirmRequest struct {
	Email string `json:"email" binding:"required"`
	Code  int64  `json:"code"  binding:"required"`
}

type RegMailRequest struct {
	Email string `json:"email" binding:"required"`
}

type TokenRequest struct {
	Email    string `json:"email" db:"email" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
}
