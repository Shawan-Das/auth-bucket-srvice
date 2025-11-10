package model

type ValidateAuthorizationOutput struct {
	Message   string      `json:"message"`
	IsSuccess bool        `json:"isSuccess"`
	Payload   interface{} `json:"errorMessage,omitempty"`
}

type CreateTokenOutput struct {
	Message      string `json:"message"`
	IsSuccess    bool   `json:"isSuccess"`
	Token        string `json:"token,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type ValidateTokenOutput struct {
	Message      string `json:"message"`
	IsSuccess    bool   `json:"isSuccess"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}