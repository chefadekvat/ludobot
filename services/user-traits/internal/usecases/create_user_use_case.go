package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"user-traits/gen/sql"
	"user-traits/internal/di"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists error = errors.New("user exists")
)

type UserCreationUseCase struct {
	ctx      context.Context
	dbconn   *pgx.Conn
	poolConn *pgxpool.Conn
}

func NewUserCreationUseCase(dependencies *di.Dependencies) (*UserCreationUseCase, error) {
	conn, err := dependencies.DbPoolPtr.Acquire(dependencies.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire database connection: %w", err)
	}

	return &UserCreationUseCase{
		ctx:      dependencies.Context,
		dbconn:   conn.Conn(),
		poolConn: conn,
	}, nil
}

func (uc *UserCreationUseCase) Close() {
	if uc.poolConn != nil {
		uc.poolConn.Release()
	}
}

func (uc *UserCreationUseCase) CreateUser(id int64, balance int64) error {
	tag, err := uc.dbconn.Exec(uc.ctx, sql.AddUser, id, balance)
	if err != nil {
		slog.Error(fmt.Sprintf("error on user %d creation: %s", id, err.Error()))
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user %d creation error: %w", id, ErrUserExists)
	}

	return nil
}
