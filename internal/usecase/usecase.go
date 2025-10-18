package usecase

import (
	"auth_service/internal/adapter/driven/dbstore"
	"auth_service/internal/config"
	"auth_service/internal/port/usecase"
	"auth_service/internal/usecase/authenticator"
	usercreator "auth_service/internal/usecase/user_creator"
)

type UseCases struct {
	UserCreator   usecase.UserCreator
	Authenticator usecase.Authenticate
}

func New(cfg config.Config, store *dbstore.DBStore) *UseCases {
	return &UseCases{
		UserCreator:   usercreator.New(&cfg, store.UserStorage),
		Authenticator: authenticator.New(&cfg, store.UserStorage),
	}
}
