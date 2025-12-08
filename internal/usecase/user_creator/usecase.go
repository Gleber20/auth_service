package user_creator

import (
	"auth_service/internal/adapter/driven/amqp"
	"auth_service/internal/config"
	"auth_service/internal/domain"
	"auth_service/internal/errs"
	"auth_service/internal/port/driven"
	"auth_service/utils"
	"context"
	"errors"
	"fmt"
)

type UseCase struct {
	cfg         *config.Config
	userStorage driven.UserStorage
	producer    driven.AMQPProducer
}

func New(cfg *config.Config, userStorage driven.UserStorage, producer driven.AMQPProducer) *UseCase {
	return &UseCase{
		cfg:         cfg,
		userStorage: userStorage,
		producer:    producer,
	}
}

func (u *UseCase) CreateUser(ctx context.Context, user domain.User) (err error) {
	// Проверить существует ли пользователь с таким username'ом в бд
	_, err = u.userStorage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotfound) {
			return err
		}
	} else {
		return errs.ErrUsernameAlreadyExists
	}

	// За хэшировать пароль
	user.Password, err = utils.GenerateHash(user.Password)
	if err != nil {
		return err
	}

	// Добавить пользователя в бд
	if err = u.userStorage.CreateUser(ctx, user); err != nil {
		return err
	}

	msg := amqp.Message{
		Recipient: user.Email,
		Subject:   "Добро пожаловать!",
		Body:      fmt.Sprintf("Привет, %s! Спасибо за регистрацию!", user.FullName),
	}
	err = u.producer.Publish(ctx, msg)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
