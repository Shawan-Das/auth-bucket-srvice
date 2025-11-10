package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "encoding/json"
)

// // CreateToken function creates a token that takes userID - if ok return token, otherwise error
// func CreateToken(id int, username string, email string, fullname string) (string, error) {
// 	var err error
// 	atClaims := jwt.MapClaims{}
// 	atClaims["authorized"] = true
// 	atClaims["id"] = id
// 	atClaims["username"] = username
// 	atClaims["email"] = email
// 	atClaims["fullname"] = fullname
// 	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
// 	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
// 	token, err := at.SignedString([]byte(ViperReturnStringConfigVariableFromLocalConfigJSON("schema_name")))
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }

// GenerateJWTAuthToken -
func GenerateJWTAuthToken(secretKey string, username string) (string, time.Time, error) {
	// func GenerateJWTAuthToken(secretKey string, payload []byte) (string, time.Time, error) {
	type Claims struct {
		// Payload []byte `json:"payload"`
		// Authorized bool   `json:"authorized"`
		// Email      string `json:"email"`
		// Fullname   string `json:"fullname"`
		// ID         int    `json:"id"`
		Username string `json:"username"`
		jwt.StandardClaims
	}

	jwtKey := []byte(secretKey)

	// Set expiration to 5 minutes from now (Maybe lesser?)
	expirationTime := time.Now().Add(5 * time.Minute)
	// create the payload
	claims := &Claims{
		// Payload: payload,
		// Authorized:
		// Email:
		// Fullname:
		// ID:
		Username:       username,
		StandardClaims: jwt.StandardClaims{
			// ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error generating JWT token : %+v", err)
	}
	return tokenString, expirationTime, nil
}

// VerifyJWTToken -
func VerifyJWTToken(secretKey string, tokenString string, username string) (*jwt.Token, string, error) {
	// func VerifyJWTToken(secretKey string, tokenString string, username string) (*jwt.Token, []byte, error) {
	type Claims struct {
		// Payload []byte `json:"payload"`
		// Authorized bool   `json:"authorized"`
		// Email      string `json:"email"`
		// Fullname   string `json:"fullname"`
		// ID         int    `json:"id"`
		Username string `json:"username"`
		jwt.StandardClaims
	}

	jwtKey := []byte(secretKey)
	// var payload []byte
	// claims := &Claims{}
	claims := &Claims{
		// Payload: payload,
		// Authorized:
		// Email:
		// Fullname:
		// ID:
		Username:       username,
		StandardClaims: jwt.StandardClaims{
			// ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtKey, nil
		})
	if err != nil {
		return nil, username, err
	}
	fmt.Println("token is verified")
	return token, claims.Username, nil
}

//  the above GenerateJWTAuthToken and VerifyJWTToken code taken from :
// https://stackoverflow.com/questions/57611134/converting-a-jwt-payload-back-to-a-struct

// // ValidateTokenWithClaims -
// func ValidateTokenWithClaims(tokenString string) {
// 	// type customClaims struct {
// 	// 	Authorized bool   `json:"authorized"`
// 	// 	Email      string `json:"email"`
// 	// 	Fullname   string `json:"fullname"`
// 	// 	ID         int    `json:"id"`
// 	// 	Username   string `json:"username"`
// 	// }
// 	// token, err := jwt.ParseWithClaims(
// 	// 	tokenString,
// 	// 	jwt.Claims
// 	// 	&customClaims{},
// 	// 	func(token *jwt.Token) (interface{}, error) {
// 	// 		return []byte(ViperReturnStringConfigVariableFromLocalConfigJSON("access_secret")), nil
// 	// 	},
// 	// )
// 	// if err != nil {
// 	// 	panic(err)
// 	// } else {
// 	// 	fmt.Println(token)
// 	// }
// }

// // VerifyJWTToken -
// func VerifyJWTToken(tokenString string) (*jwt.Token, []byte, error) {
// 	var jwtKey = []byte(ViperReturnStringConfigVariableFromLocalConfigJSON("access_secret"))
// 	var payload []byte
// 	// claims := &Claims{}

// 	type claims struct{

// 	}

// 	token, err := jwt.ParseWithClaims(tokenString, claims,
// 		func(token *jwt.Token) (interface{}, error) {

// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 			}

// 			return jwtKey, nil
// 		})

// 	if err != nil {
// 		return nil, payload, err
// 	}

// 	return token, claims.Payload, nil
// }
