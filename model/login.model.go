package model

// Login struct for db table - login
type Login struct {
	Username     string `json:"username"`
	Userpassword string `json:"userpassword"`
	Useremail    string `json:"useremail"`
	Token        string `json:"token"`
}

type Logins struct {
	Id            int    `json:"id"`
	Fullname      string `json:"fullname"`
	Username      string `json:"username"`
	Hash_password string `json:"hash_password"`
	Otp_hash      string `json:"otp_hash"`
	Pass_flag     int    `json:"pass_flag"`
	Token_value   string `json:"token_value"`
	Phoneno       string `json:"phoneno"`
	Otp_verified  bool   `json:"otp_verified"`
	Is_access     bool   `json:"is_access"`
	User_lock     bool   `json:"user_lock"`
}
