package customErrors

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var (
	ErrNotFound     = pgx.ErrNoRows
	ErrConflict     = errors.New("conflict")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
)

func GetPgErrCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func HandlePgError(err error, msg string) error {
	if err == ErrNotFound {
		return ErrNotFound
	}

	code := GetPgErrCode(err)
	switch code {
	case UniqueViolation:
		return errors.Wrap(ErrConflict, "email already exists")
	default:
		return errors.Wrap(err, msg)
	}
}
