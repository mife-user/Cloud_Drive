package response

import "drive/internal/domain"

type UserRegRPS struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

type UserLoginRPS struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

// ToDTUserReg 转换为用户注册响应DTO
func ToDTUserReg(user *domain.User) *UserRegRPS {
	return &UserRegRPS{
		UserID: user.ID,
		Role:   user.Role,
	}
}

// ToDTUserLogin 转换为用户登录响应DTO
func ToDTUserLogin(user *domain.User, token string) *UserLoginRPS {
	return &UserLoginRPS{
		UserID:   user.ID,
		UserName: user.UserName,
		Role:     user.Role,
		Token:    token,
	}
}
