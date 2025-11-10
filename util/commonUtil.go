package util

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/alexedwards/argon2id"
)

var ConfigFileName string

func TrimString(s string) string {
	input := strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(s, " "))
	return input
}

func IsValidPass(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func HashString(input string) string {
	hash, err := argon2id.CreateHash(input, argon2id.DefaultParams)
	if err != nil {
		return ""
	}
	return hash
}

func DoStringMatch(password string, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false
	}
	return match
}

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	otpint := rand.Intn(max-min+1) + min
	otp := strconv.Itoa(otpint)

	return otp
}

func Get3SubstringsBasedOnCapital(input string) string {
	var result strings.Builder
	currentSubstring := ""

	for _, char := range input {
		if unicode.IsUpper(char) {
			if len(currentSubstring) < 3 {
				currentSubstring += string(char)
			} else {
				result.WriteString(currentSubstring)
				currentSubstring = string(char)
			}
		} else {
			if len(currentSubstring) > 0 {
				result.WriteString(currentSubstring)
				currentSubstring = ""
			}
		}
	}

	if len(currentSubstring) > 0 {
		result.WriteString(currentSubstring)
	}

	return result.String()
}

func Get4SubstringsBasedOnCapital(input string) string {
	var result strings.Builder
	currentSubstring := ""

	for _, char := range input {
		if unicode.IsUpper(char) {
			if len(currentSubstring) < 4 {
				currentSubstring += string(char)
			} else {
				result.WriteString(currentSubstring)
				currentSubstring = string(char)
			}
		} else {
			if len(currentSubstring) > 0 {
				result.WriteString(currentSubstring)
				currentSubstring = ""
			}
		}
	}

	if len(currentSubstring) > 0 {
		result.WriteString(currentSubstring)
	}

	return result.String()
}

func LoadConfig(filePath string) (map[string]interface{}, error) {
	config := make(map[string]interface{})

	configFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ContainsString(arr []string, value string) bool {
	for _, v := range arr {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

func ContainsSubstring(arr []string, value string) bool {
	for _, v := range arr {
		if strings.Contains(value, v) {
			return true
		}
	}
	return false
}


