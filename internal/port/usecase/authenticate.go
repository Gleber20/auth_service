package usecase

import (
	"auth_service/internal/domain"
	"context"
)

type Authenticate interface {
	Authenticate(ctx context.Context, user domain.User) (int, error)
}
