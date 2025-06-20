package dtos

import "github.com/yelaco/gchess-server/internal/pkg/domains/models/entities"

type UserCreateRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,validpassword"`
}

type UserUpdateRequest struct {
	ID        uint64
	Email     string `json:"email" binding:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password" binding:"gte=8,lte=32"`
}

type UserResponse struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func UserResponseFromEntity(user entities.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
