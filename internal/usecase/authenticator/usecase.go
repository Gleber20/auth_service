package authenticator

import (
	"auth_service/internal/config"
	"auth_service/internal/domain"
	"auth_service/internal/errs"
	"auth_service/internal/port/driven"
	"auth_service/utils"
	"context"
	"errors"
)

type UseCase struct {
	cfg         *config.Config
	userStorage driven.UserStorage
}

func New(cfg *config.Config, userStorage driven.UserStorage) *UseCase {
	return &UseCase{
		cfg:         cfg,
		userStorage: userStorage,
	}
}

func (u *UseCase) Authenticate(ctx context.Context, user domain.User) (int, error) {
	// проверить существует ли пользователь с таким username
	userFromDB, err := u.userStorage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotfound) {
			return 0, errs.ErrUserNotFound
		}

		return 0, err
	}

	// за хэшировать пароль, который получили от пользователя
	user.Password, err = utils.GenerateHash(user.Password)
	if err != nil {
		return 0, err
	}

	// проверить правильно ли он указал пароль
	if userFromDB.Password != user.Password {
		return 0, errs.ErrIncorrectUsernameOrPassword
	}

	return userFromDB.ID, nil
}
