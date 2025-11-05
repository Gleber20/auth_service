package dbstore

import (
	"auth_service/internal/domain"
	"context"
	"os"
	_ "os/user"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

type User struct {
	ID        int       `db:"id"`
	FullName  string    `db:"full_name"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		FullName:  u.FullName,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) FromDomain(d domain.User) {
	u.ID = d.ID
	u.FullName = d.FullName
	u.Username = d.Username
	u.Email = d.Email
	u.Password = d.Password
	u.UpdatedAt = d.UpdatedAt
	u.CreatedAt = d.CreatedAt
}

func (u *UserStorage) CreateUser(ctx context.Context, user domain.User) (err error) {
	var dbUser User
	dbUser.FromDomain(user)

	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "CreateUser").Logger()
	_, err = u.db.ExecContext(ctx, `INSERT INTO users (full_name, username, email, password)
					VALUES ($1, $2, $3, $4)`,
		dbUser.FullName,
		dbUser.Username,
		dbUser.Email,
		dbUser.Password)
	if err != nil {
		logger.Err(err).Msg("error inserting user")
		return u.translateError(err)
	}

	return nil
}

func (u *UserStorage) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "GetUserByID").Logger()

	var dbUser User
	if err := u.db.GetContext(ctx, &dbUser, `
		SELECT id, full_name, username, email, password, created_at, updated_at 
		FROM users
		WHERE id = $1`, id); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, u.translateError(err)
	}
	return *dbUser.ToDomain(), nil
}

func (u *UserStorage) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "GetUserByUsername").Logger()

	var dbUser User
	if err := u.db.GetContext(ctx, &dbUser, `
		SELECT id, full_name, username, email, password, created_at, updated_at 
		FROM users
		WHERE username = $1`, username); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, u.translateError(err)
	}

	return *dbUser.ToDomain(), nil
}
