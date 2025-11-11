package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	auth "github.com/satcomit/hrms/internal/dbmodel/db_query"

	"github.com/satcomit/hrms/internal/model"
	"github.com/satcomit/hrms/internal/util"

	"github.com/sirupsen/logrus"
)

var _asLogger = logrus.New()

// AuthenticationRESTService provides authentication related rest services
type AuthenticationRESTService struct {
	dbConn        *util.DBConnectionWrapper
	jwtSigningKey []byte
	bypassAuth    map[string]bool
}

// NewAuthenticationRESTService returns a new initialized version of the service
func NewAuthenticationRESTService(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) *AuthenticationRESTService {
	service := new(AuthenticationRESTService)
	if err := service.Init(config, dbConnection, verbose); err != nil {
		_asLogger.Errorf("Unable to initialize service instance %v", err)
		return nil
	}
	return service
}

// Init initializes the service instance
func (s *AuthenticationRESTService) Init(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) error {
	if verbose {
		_asLogger.SetLevel(logrus.DebugLevel)
	}
	if dbConnection == nil {
		return fmt.Errorf("null DB Util reference passed")
	}
	s.dbConn = dbConnection
	var conf model.AuthServiceConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		_asLogger.Error("Unable to parse config json file ", err)
		return err
	}
	if conf.JWTKey != nil && len(*conf.JWTKey) > 0 {
		s.jwtSigningKey = []byte(*conf.JWTKey)
	}
	s.bypassAuth = make(map[string]bool)
	s.bypassAuth["/"] = true
	if conf.BypassAuth != nil && len(conf.BypassAuth) > 0 {
		for _, url := range conf.BypassAuth {
			s.bypassAuth[url] = true
		}
	}
	_asLogger.Infof("Successfully initialized AuthenticationRESTService")
	return nil
}

// AddRouters add api end points specific to this service
func (s *AuthenticationRESTService) AddRouters(apiBase string, router *gin.Engine) {
	router.POST("/api/auth/create", func(c *gin.Context) {
		resp := s.createUser(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.POST("/api/auth/login", func(c *gin.Context) {
		resp := s.validateLogin(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.POST("/api/auth/resetpwd", func(c *gin.Context) {
		resp := s.resetPassword(c)
		c.JSON(resp.StatusCode, resp)
	})

	router.GET("/api/auth/users", func(c *gin.Context) {
		resp := s.getAllUsers(c)
		c.JSON(resp.StatusCode, resp)
	})
}

// /api/auth/create - create user
func (s *AuthenticationRESTService) createUser(c *gin.Context) APIResponse {
	var input model.CreateUserInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	// Validate required fields
	if input.Email == "" || input.Password == "" || input.Phone == "" {
		return BuildResponse400("Email, password, and phone are required")
	}

	// Check if user already exists
	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	_, err := qtx.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return BuildResponse400("User with this email already exists")
	}

	// Set default values
	userName := input.UserName
	if userName == "" {
		userName = input.Email // Use email as default username
	}
	role := input.Role
	if role == "" {
		role = "USER" // Default role
	}

	// Hash password
	hashedPassword := s.getHashOf(input.Password)

	// Create user
	createParams := auth.CreateUserParams{
		UserName: userName,
		Email:    input.Email,
		Phone:    input.Phone,
		Pass:     hashedPassword,
		Role:     role,
	}

	err = qtx.CreateUser(ctx, createParams)
	if err != nil {
		_asLogger.Errorf("Error creating user: %v", err)
		return BuildResponse500("Failed to create user", err.Error())
	}

	return BuildResponse200("User created successfully", nil)
}

// /api/auth/login - login (supports username, email, or phone)
func (s *AuthenticationRESTService) validateLogin(c *gin.Context) APIResponse {
	var input model.LoginInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	if input.Login == "" || input.Password == "" {
		return BuildResponse400("Login identifier (username/email/phone) and password are required")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Try to find user by username, email, or phone
	user, err := qtx.GetUserByLogin(ctx, input.Login)
	if err != nil {
		_asLogger.Errorf("Error getting user: %v", err)
		return BuildResponse404("Invalid login credentials or password", false)
	}

	// Check password
	hashedPassword := s.getHashOf(input.Password)
	if user.Pass != hashedPassword {
		return BuildResponse404("Invalid login credentials or password", false)
	}

	// Check if password is valid
	if !user.PssValid {
		return BuildResponse400("Password is not valid. Please reset your password")
	}

	// Create JWT token
	jwtToken := s.createJWTToken(user.UserID, user.Email, user.UserName)

	response := BuildResponse200("Login successful", map[string]interface{}{
		"user_id":   user.UserID,
		"user_name": user.UserName,
		"email":     user.Email,
		"phone":     user.Phone,
		"role":      user.Role,
	})
	response.Token = &jwtToken

	return response
}

// /api/auth/resetpwd - reset password
func (s *AuthenticationRESTService) resetPassword(c *gin.Context) APIResponse {
	var input model.AuthDataInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	if input.Email == "" || input.NewPassword == "" {
		return BuildResponse400("Email and new password are required")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Check if user exists
	_, err := qtx.GetUserByEmail(ctx, input.Email)
	if err != nil {
		_asLogger.Errorf("Error getting user: %v", err)
		return BuildResponse404("User not found", false)
	}

	// Hash new password
	hashedPassword := s.getHashOf(input.NewPassword)

	// Update password
	updateParams := auth.UpdatePasswordParams{
		Pass:     hashedPassword,
		PssValid: true,
		Email:    input.Email,
	}

	err = qtx.UpdatePassword(ctx, updateParams)
	if err != nil {
		_asLogger.Errorf("Error updating password: %v", err)
		return BuildResponse500("Failed to reset password", err.Error())
	}

	return BuildResponse200("Password reset successfully", nil)
}

// /api/auth/users - get all users (requires authentication)
func (s *AuthenticationRESTService) getAllUsers(c *gin.Context) APIResponse {
	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	users, err := qtx.GetAllUsers(ctx)
	if err != nil {
		_asLogger.Errorf("Error getting users: %v", err)
		return BuildResponse500("Failed to retrieve users", err.Error())
	}

	// Transform to response format
	userList := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		userList = append(userList, map[string]interface{}{
			"code":  user.UserID,
			"name":  user.UserName,
			"email": user.Email,
		})
	}

	return BuildResponse200("Users retrieved successfully", userList)
}

func (s *AuthenticationRESTService) getHashOf(password string) string {
	shaBytes := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", shaBytes)
}

func (s *AuthenticationRESTService) createJWTToken(userID int32, email, userName string) string {
	if s.jwtSigningKey == nil {
		return ""
	}
	claim := model.AuthorizationClaims{
		UserID:   userID,
		Email:    email,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			Issuer:    "Auth Service",
			Id:        fmt.Sprintf("%d", userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(s.jwtSigningKey)
	if err != nil {
		_asLogger.Error("Error in generating token", err)
		return ""
	}
	_asLogger.Infof("Generated token for user %s", email)

	return tokenStr
}

func (s *AuthenticationRESTService) checkAuth(c *gin.Context) bool {
	url := c.Request.URL
	uri := url.RequestURI()

	// Allow bypass URLs from config
	if _, isFound := s.bypassAuth[uri]; isFound {
		return true
	}

	// If JWT key is not set, allow all
	if s.jwtSigningKey == nil {
		return true
	}

	// Check for JWT token in Authorization header
	authHeader := c.Request.Header.Get("Authorization")
	if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSigningKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	// Token is valid, allow request
	return true
}
