package usecase

import (
	"auth_service/internal/domain"
	"context"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user domain.User) (err error)
}
