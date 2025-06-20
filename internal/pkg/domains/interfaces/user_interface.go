package interfaces

import (
	"github.com/yelaco/gchess-server/internal/pkg/domains/models/dtos"
	"github.com/yelaco/gchess-server/internal/pkg/domains/models/entities"
)

type UserUsecase interface {
	CreateUser(dtos.UserCreateRequest) (entities.User, error)
	GetUserByID(id uint64) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	UpdateUser(dtos.UserUpdateRequest) (entities.User, error)
	DeleteUserByID(id uint64) error
}
