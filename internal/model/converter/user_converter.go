package converter

import (
	"go-rest-scaffold/internal/entity"
	"go-rest-scaffold/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToTokenResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		Token:        user.Token,
		RefreshToken: user.RefreshToken,
	}
}
