package utils

type AddUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phonenumber,omitempty"`
	IsAdmin     bool   `json:"isadmin"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
