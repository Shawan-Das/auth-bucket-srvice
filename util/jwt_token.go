package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type JWTClaims struct {
	UserID   int32  `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"`
	jwt.StandardClaims
}

type TokenPair struct {
	AccessToken         string
	RefreshToken        string
	AccessTokenExpiry   time.Time
	RefreshTokenExpiry  time.Time
	AccessTokenDuration time.Duration
}

func GetJWTSecret() string {
	secret := viper.GetString("access_secret")
	if secret == "" {
		secret = "default-secret-key-please-change-in-production"
	}
	return secret
}

func GetAccessTokenDuration() time.Duration {
	duration := viper.GetInt("access_token_duration_minutes")
	if duration <= 0 {
		duration = 15
	}
	return time.Duration(duration) * time.Minute
}

func GetRefreshTokenDuration() time.Duration {
	duration := viper.GetInt("refresh_token_duration_hours")
	if duration <= 0 {
		duration = 168
	}
	return time.Duration(duration) * time.Hour
}

func GenerateTokenPair(userID int32, email, username, role string) (*TokenPair, error) {
	secret := GetJWTSecret()
	accessDuration := GetAccessTokenDuration()
	refreshDuration := GetRefreshTokenDuration()

	accessExpiry := time.Now().Add(accessDuration)
	refreshExpiry := time.Now().Add(refreshDuration)

	accessClaims := &JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		Type:     string(AccessToken),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessTokenObj.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshClaims := &JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		Type:     string(RefreshToken),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshTokenObj.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:         accessTokenStr,
		RefreshToken:        refreshTokenStr,
		AccessTokenExpiry:   accessExpiry,
		RefreshTokenExpiry:  refreshExpiry,
		AccessTokenDuration: accessDuration,
	}, nil
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	secret := GetJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != string(AccessToken) {
		return nil, fmt.Errorf("invalid token type: expected access token")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != string(RefreshToken) {
		return nil, fmt.Errorf("invalid token type: expected refresh token")
	}

	return claims, nil
}
