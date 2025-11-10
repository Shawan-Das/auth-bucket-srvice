package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tools/common/model"
	"github.com/tools/common/util"
)

type Authorization struct {
	tokenRepository *TokenRepository
}

func (authorizationrepository *Authorization) ValidateAuthorization_V2(c *gin.Context) model.ValidateAuthorizationOutput {
	url := c.Request.URL
	uri := url.RequestURI()

	authorizationUrls := viper.GetViper().GetStringMapString("bypassAuth_V2")

	equalFoldUrls, equalFoldOk := authorizationUrls["reqBody"]
	containsUrls, containsOk := authorizationUrls["params"]

	if !equalFoldOk || !containsOk {
		return model.ValidateAuthorizationOutput{
			Message:   "Invalid authorizationUrls format in config.",
			IsSuccess: false,
		}
	}

	equalFoldStrUrls := strings.Split(equalFoldUrls, ",")
	containsStrUrls := strings.Split(containsUrls, ",")

	if util.ContainsString(equalFoldStrUrls, uri) || util.ContainsSubstring(containsStrUrls, uri) {
		return model.ValidateAuthorizationOutput{
			Message:   "Authorization successful.",
			IsSuccess: true,
		}
	}

	authHeader := c.Request.Header.Get("Authorization")
	if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer") {
		return model.ValidateAuthorizationOutput{
			Message:   "Invalid authotization header. Please login again.",
			IsSuccess: false,
		}
	}

	authArray := strings.Split(authHeader, " ")
	if len(authArray) != 2 {
		return model.ValidateAuthorizationOutput{
			Message:   "Unable to parse authorization token. Please login again.",
			IsSuccess: false,
		}
	}

	validateTokenOutput := authorizationrepository.tokenRepository.ValidateToken(authArray[1])
	if !validateTokenOutput.IsSuccess {
		return model.ValidateAuthorizationOutput{
			Message:   "Authorization failed: Token is invalid. Please login again.",
			IsSuccess: false,
			Payload:   validateTokenOutput,
		}
	}

	return model.ValidateAuthorizationOutput{
		Message:   "Authorization successful. Token is valid.",
		IsSuccess: true,
	}
}
