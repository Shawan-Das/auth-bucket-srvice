package service

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/tools/common/model"
	"github.com/tools/common/util"
)

type TokenRepository struct{}

func (tokenRepository *TokenRepository) CreateToken(username string) model.CreateTokenOutput {
	if username == "" {
		return model.CreateTokenOutput{
			Message:   "Token creation failed. Username is empty",
			IsSuccess: false,
		}
	}
	mySigningKey := []byte(util.ViperReturnStringConfigVariableFromLocalConfigJSON("ACCESS_SECRET"))
	type MyCustomClaims struct {
		Preferred_username string `json:"preferred_username"`
		jwt.StandardClaims
	}
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			Audience:  username,
			ExpiresAt: time.Now().Add(30 * time.Minute * 60).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    "APIP_SERVICE",
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return model.CreateTokenOutput{
			Message:      "Token creation failed. Something went wrong. Please login again.",
			IsSuccess:    false,
			Token:        ss,
			ErrorMessage: err.Error(),
		}
	}
	return model.CreateTokenOutput{
		Message:   "Token creation successful",
		IsSuccess: true,
		Token:     ss,
	}
}

func (tokenRepository *TokenRepository) ValidateToken(tokenString string) model.ValidateTokenOutput {
	if tokenString == "" {
		return model.ValidateTokenOutput{
			Message:   "Token is empty. Please login again",
			IsSuccess: false,
		}
	}
	type MyCustomClaims struct {
		Preferred_username string `json:"preferred_username"`
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(util.ViperReturnStringConfigVariableFromLocalConfigJSON("ACCESS_SECRET")), nil
	})
	if err != nil {
		return model.ValidateTokenOutput{
			Message:      "Token validation failed. Please login again.",
			IsSuccess:    false,
			ErrorMessage: err.Error(),
		}
	}
	if _, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return model.ValidateTokenOutput{
			Message:   "Token validated successfully",
			IsSuccess: true,
		}
	}
	return model.ValidateTokenOutput{
		Message:   "Token validation failed. Something went wrong. Please login again",
		IsSuccess: false,
	}
}
