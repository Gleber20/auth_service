package dbstore

import (
	"auth_service/internal/errs"
	"database/sql"
	"errors"
)

func (u *UserStorage) translateError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotfound
	default:
		return err
	}
}
