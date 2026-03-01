package response

import "drive/internal/domain"

func ToDTUser(user *domain.User) *UserRegResponse {
	return &UserRegResponse{
		UserID: user.ID,
		Role:   user.Role,
	}
}
