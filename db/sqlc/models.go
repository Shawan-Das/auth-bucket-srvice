package db

import (
	"time"
)

type User struct {
	UserID          int32      `json:"user_id"`
	UserName        string     `json:"user_name"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Pass            string     `json:"pass"`
	PssValid        bool       `json:"pss_valid"`
	Otp             string     `json:"otp"`
	OtpValid        bool       `json:"otp_valid"`
	OtpExp          *time.Time `json:"otp_exp"`
	Role            string     `json:"role"`
	RefreshToken    string     `json:"refresh_token"`
	RefreshTokenExp *time.Time `json:"refresh_token_exp"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
