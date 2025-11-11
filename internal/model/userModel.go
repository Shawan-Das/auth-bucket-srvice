package model

type LoginInput struct {
	Login    string `json:"login"` // Can be username, email, or phone
	Password string `json:"pwd"`
}

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	UserName string `json:"userName,omitempty"`
	Role     string `json:"role,omitempty"`
}
