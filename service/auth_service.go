package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/tools/common/dto"
	"github.com/tools/common/repository"
	"github.com/tools/common/util"
)

var _authLogger = logrus.New()

type AuthService struct {
	userRepo *repository.UserRepository
	verbose  bool
}

func NewAuthService(conn *pgx.Conn, verbose bool) *AuthService {
	if verbose {
		_authLogger.SetLevel(logrus.DebugLevel)
	}

	return &AuthService{
		userRepo: repository.NewUserRepository(conn, verbose),
		verbose:  verbose,
	}
}

func (s *AuthService) AddRouters(router *gin.Engine) {
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", s.Register)
		authGroup.POST("/login", s.Login)
		authGroup.POST("/refresh", s.RefreshToken)
		authGroup.POST("/logout", s.Logout)
	}
}

func (s *AuthService) Register(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_authLogger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, buildResponse(false, "Invalid request body", nil))
		return
	}

	if !util.IsEmailValid(req.Email) {
		c.JSON(http.StatusBadRequest, buildResponse(false, "Invalid email format", nil))
		return
	}

	if !util.IsValidPass(req.Password) {
		c.JSON(http.StatusBadRequest, buildResponse(false, "Password must be at least 8 characters and contain uppercase, lowercase, number, and special character", nil))
		return
	}

	ctx := context.Background()

	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, buildResponse(false, "User with this email already exists", nil))
		return
	}

	user, err := s.userRepo.CreateUser(ctx, req.UserName, req.Email, req.Phone, req.Password, "ADMIN")
	if err != nil {
		_authLogger.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, buildResponse(false, "Failed to create user", nil))
		return
	}

	response := dto.UserResponse{
		UserID:    user.UserID,
		UserName:  user.UserName,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, buildResponse(true, "User registered successfully", response))
}

func (s *AuthService) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_authLogger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, buildResponse(false, "Invalid request body", nil))
		return
	}

	ctx := context.Background()

	user, err := s.userRepo.ValidateUserPassword(ctx, req.Email, req.Password)
	if err != nil {
		_authLogger.Errorf("Authentication failed: %v", err)
		c.JSON(http.StatusUnauthorized, buildResponse(false, "Invalid email or password", nil))
		return
	}

	if user.Role != "ADMIN" {
		c.JSON(http.StatusForbidden, buildResponse(false, "Access denied. Admin role required", nil))
		return
	}

	tokenPair, err := util.GenerateTokenPair(user.UserID, user.Email, user.UserName, user.Role)
	if err != nil {
		_authLogger.Errorf("Failed to generate tokens: %v", err)
		c.JSON(http.StatusInternalServerError, buildResponse(false, "Failed to generate authentication tokens", nil))
		return
	}

	err = s.userRepo.UpdateRefreshToken(ctx, user.UserID, tokenPair.RefreshToken, tokenPair.RefreshTokenExpiry)
	if err != nil {
		_authLogger.Errorf("Failed to save refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, buildResponse(false, "Failed to complete authentication", nil))
		return
	}

	response := dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		UserID:       user.UserID,
		UserName:     user.UserName,
		Email:        user.Email,
		Role:         user.Role,
		ExpiresIn:    int64(tokenPair.AccessTokenDuration.Seconds()),
	}

	_authLogger.Infof("User logged in successfully: %s", user.Email)
	c.JSON(http.StatusOK, buildResponse(true, "Login successful", response))
}

func (s *AuthService) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_authLogger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, buildResponse(false, "Invalid request body", nil))
		return
	}

	claims, err := util.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		_authLogger.Errorf("Invalid refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, buildResponse(false, "Invalid or expired refresh token", nil))
		return
	}

	ctx := context.Background()

	user, err := s.userRepo.GetUserByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		_authLogger.Errorf("Refresh token not found in database: %v", err)
		c.JSON(http.StatusUnauthorized, buildResponse(false, "Invalid or expired refresh token", nil))
		return
	}

	if user.UserID != claims.UserID {
		_authLogger.Errorf("User ID mismatch in refresh token")
		c.JSON(http.StatusUnauthorized, buildResponse(false, "Invalid refresh token", nil))
		return
	}

	tokenPair, err := util.GenerateTokenPair(user.UserID, user.Email, user.UserName, user.Role)
	if err != nil {
		_authLogger.Errorf("Failed to generate new tokens: %v", err)
		c.JSON(http.StatusInternalServerError, buildResponse(false, "Failed to generate new tokens", nil))
		return
	}

	err = s.userRepo.UpdateRefreshToken(ctx, user.UserID, tokenPair.RefreshToken, tokenPair.RefreshTokenExpiry)
	if err != nil {
		_authLogger.Errorf("Failed to update refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, buildResponse(false, "Failed to refresh token", nil))
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    int64(tokenPair.AccessTokenDuration.Seconds()),
	}

	_authLogger.Infof("Token refreshed successfully for user: %s", user.Email)
	c.JSON(http.StatusOK, buildResponse(true, "Token refreshed successfully", response))
}

func (s *AuthService) Logout(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, buildResponse(false, "Unauthorized", nil))
		return
	}

	userID, ok := userIDInterface.(int32)
	if !ok {
		c.JSON(http.StatusInternalServerError, buildResponse(false, "Invalid user ID format", nil))
		return
	}

	ctx := context.Background()

	err := s.userRepo.InvalidateRefreshToken(ctx, userID)
	if err != nil {
		_authLogger.Errorf("Failed to invalidate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, buildResponse(false, "Failed to logout", nil))
		return
	}

	_authLogger.Infof("User logged out successfully: %d", userID)
	c.JSON(http.StatusOK, buildResponse(true, "Logout successful", nil))
}

type Authorization struct{}

func (a *Authorization) ValidateAuthorization(c *gin.Context) APIResponse {
	path := c.Request.URL.Path

	bypassPaths := []string{"/api/auth/register", "/api/auth/login", "/api/auth/refresh", "/", "/swagger"}
	for _, bypassPath := range bypassPaths {
		if path == bypassPath || (len(path) >= len(bypassPath) && path[:len(bypassPath)] == bypassPath) {
			return APIResponse{IsSuccess: true}
		}
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return APIResponse{
			IsSuccess: false,
			Message:   "Authorization header required",
		}
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return APIResponse{
			IsSuccess: false,
			Message:   "Invalid authorization format",
		}
	}

	tokenString := authHeader[7:]

	claims, err := util.ValidateAccessToken(tokenString)
	if err != nil {
		return APIResponse{
			IsSuccess: false,
			Message:   fmt.Sprintf("Invalid or expired token: %v", err),
		}
	}

	c.Set("userID", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)

	return APIResponse{IsSuccess: true}
}
