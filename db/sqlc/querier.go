package db

import (
	"context"
	"time"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, userID int32) (User, error)
	UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) error
	GetUserByRefreshToken(ctx context.Context, refreshToken string) (User, error)
	InvalidateRefreshToken(ctx context.Context, userID int32) error
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	GetUserByEmailWithRole(ctx context.Context, arg GetUserByEmailWithRoleParams) (User, error)
}

var _ Querier = (*Queries)(nil)
