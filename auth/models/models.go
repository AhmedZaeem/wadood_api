package models

import "time"

type User struct {
	ID          int               `json:"id"`
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	PhoneNumber string            `json:"phone_number"`
	Password    string            `json:"-"`
	Tokens      map[string]string `json:"tokens"`
	LastLogin   time.Time         `json:"last_login"`
}

type Device struct {
	ID    int    `json:"id"`
	IMEI  string `json:"imei"`
	Token string `json:"token"`
}

type AuthResponse struct {
	User    *User  `json:"user"`
	Token   string `json:"token"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type EditProfileRequest struct {
	Language string `json:"language"`
}
