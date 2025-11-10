package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

const createUser = `-- name: CreateUser :one
INSERT INTO common.users (
    user_name,
    email,
    phone,
    pass,
    role,
    pss_valid,
    otp_valid
) VALUES (
    $1, $2, $3, $4, $5, true, false
) RETURNING user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role, refresh_token, refresh_token_exp, created_at, updated_at
`

type CreateUserParams struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Pass     string `json:"pass"`
	Role     string `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.UserName,
		arg.Email,
		arg.Phone,
		arg.Pass,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.Email,
		&i.Phone,
		&i.Pass,
		&i.PssValid,
		&i.Otp,
		&i.OtpValid,
		&i.OtpExp,
		&i.Role,
		&i.RefreshToken,
		&i.RefreshTokenExp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role, refresh_token, refresh_token_exp, created_at, updated_at FROM common.users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.Email,
		&i.Phone,
		&i.Pass,
		&i.PssValid,
		&i.Otp,
		&i.OtpValid,
		&i.OtpExp,
		&i.Role,
		&i.RefreshToken,
		&i.RefreshTokenExp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role, refresh_token, refresh_token_exp, created_at, updated_at FROM common.users
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.Email,
		&i.Phone,
		&i.Pass,
		&i.PssValid,
		&i.Otp,
		&i.OtpValid,
		&i.OtpExp,
		&i.Role,
		&i.RefreshToken,
		&i.RefreshTokenExp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateRefreshToken = `-- name: UpdateRefreshToken :exec
UPDATE common.users
SET refresh_token = $1,
    refresh_token_exp = $2,
    updated_at = now()
WHERE user_id = $3
`

type UpdateRefreshTokenParams struct {
	RefreshToken    string     `json:"refresh_token"`
	RefreshTokenExp *time.Time `json:"refresh_token_exp"`
	UserID          int32      `json:"user_id"`
}

func (q *Queries) UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) error {
	_, err := q.db.Exec(ctx, updateRefreshToken, arg.RefreshToken, arg.RefreshTokenExp, arg.UserID)
	return err
}

const getUserByRefreshToken = `-- name: GetUserByRefreshToken :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role, refresh_token, refresh_token_exp, created_at, updated_at FROM common.users
WHERE refresh_token = $1
AND refresh_token_exp > now()
LIMIT 1
`

func (q *Queries) GetUserByRefreshToken(ctx context.Context, refreshToken string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByRefreshToken, refreshToken)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.Email,
		&i.Phone,
		&i.Pass,
		&i.PssValid,
		&i.Otp,
		&i.OtpValid,
		&i.OtpExp,
		&i.Role,
		&i.RefreshToken,
		&i.RefreshTokenExp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const invalidateRefreshToken = `-- name: InvalidateRefreshToken :exec
UPDATE common.users
SET refresh_token = '',
    refresh_token_exp = NULL,
    updated_at = now()
WHERE user_id = $1
`

func (q *Queries) InvalidateRefreshToken(ctx context.Context, userID int32) error {
	_, err := q.db.Exec(ctx, invalidateRefreshToken, userID)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE common.users
SET pass = $1,
    pss_valid = true,
    updated_at = now()
WHERE user_id = $2
`

type UpdateUserPasswordParams struct {
	Pass   string `json:"pass"`
	UserID int32  `json:"user_id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Pass, arg.UserID)
	return err
}

const getUserByEmailWithRole = `-- name: GetUserByEmailWithRole :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, otp_valid, otp_exp, role, refresh_token, refresh_token_exp, created_at, updated_at FROM common.users
WHERE email = $1 AND role = $2 LIMIT 1
`

type GetUserByEmailWithRoleParams struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (q *Queries) GetUserByEmailWithRole(ctx context.Context, arg GetUserByEmailWithRoleParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmailWithRole, arg.Email, arg.Role)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.Email,
		&i.Phone,
		&i.Pass,
		&i.PssValid,
		&i.Otp,
		&i.OtpValid,
		&i.OtpExp,
		&i.Role,
		&i.RefreshToken,
		&i.RefreshTokenExp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
