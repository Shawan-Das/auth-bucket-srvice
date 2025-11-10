package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/tools/common/db/sqlc"
	"github.com/tools/common/util"
)

var _userRepoLogger = logrus.New()

type UserRepository struct {
	conn    *pgx.Conn
	queries *db.Queries
	verbose bool
}

func NewUserRepository(conn *pgx.Conn, verbose bool) *UserRepository {
	if verbose {
		_userRepoLogger.SetLevel(logrus.DebugLevel)
	}
	return &UserRepository{
		conn:    conn,
		queries: db.New(conn),
		verbose: verbose,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, userName, email, phone, password, role string) (*db.User, error) {
	hashedPassword := util.HashString(password)
	if hashedPassword == "" {
		return nil, fmt.Errorf("failed to hash password")
	}

	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		UserName: userName,
		Email:    email,
		Phone:    phone,
		Pass:     hashedPassword,
		Role:     role,
	})

	if err != nil {
		_userRepoLogger.Errorf("Error creating user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	_userRepoLogger.Infof("User created successfully: %s", email)
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		_userRepoLogger.Errorf("Error getting user by email: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int32) (*db.User, error) {
	user, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		_userRepoLogger.Errorf("Error getting user by ID: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) ValidateUserPassword(ctx context.Context, email, password string) (*db.User, error) {
	user, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.PssValid {
		return nil, fmt.Errorf("password is not valid")
	}

	if !util.DoStringMatch(password, user.Pass) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (r *UserRepository) UpdateRefreshToken(ctx context.Context, userID int32, refreshToken string, expiry time.Time) error {
	err := r.queries.UpdateRefreshToken(ctx, db.UpdateRefreshTokenParams{
		RefreshToken:    refreshToken,
		RefreshTokenExp: &expiry,
		UserID:          userID,
	})

	if err != nil {
		_userRepoLogger.Errorf("Error updating refresh token: %v", err)
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	_userRepoLogger.Infof("Refresh token updated for user ID: %d", userID)
	return nil
}

func (r *UserRepository) GetUserByRefreshToken(ctx context.Context, refreshToken string) (*db.User, error) {
	user, err := r.queries.GetUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invalid or expired refresh token")
		}
		_userRepoLogger.Errorf("Error getting user by refresh token: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) InvalidateRefreshToken(ctx context.Context, userID int32) error {
	err := r.queries.InvalidateRefreshToken(ctx, userID)
	if err != nil {
		_userRepoLogger.Errorf("Error invalidating refresh token: %v", err)
		return fmt.Errorf("failed to invalidate refresh token: %w", err)
	}

	_userRepoLogger.Infof("Refresh token invalidated for user ID: %d", userID)
	return nil
}

func (r *UserRepository) GetUserByEmailWithRole(ctx context.Context, email, role string) (*db.User, error) {
	user, err := r.queries.GetUserByEmailWithRole(ctx, db.GetUserByEmailWithRoleParams{
		Email: email,
		Role:  role,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found with specified role")
		}
		_userRepoLogger.Errorf("Error getting user by email with role: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
