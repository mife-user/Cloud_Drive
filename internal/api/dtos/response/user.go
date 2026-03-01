package response

type UserRegResponse struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
