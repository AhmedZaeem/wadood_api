package models

import "time"

type User struct {
    ID        int               `json:"id"`
    Username  string            `json:"username"`
    Email     string            `json:"email"`
    Password  string            `json:"-"`
    Tokens    map[string]string `json:"tokens"`
    LastLogin time.Time         `json:"last_login"`
    Language  string            `json:"language"`
}

type AuthResponse struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Token   string      `json:"token,omitempty"`
}

type EditProfileRequest struct {
    Language string `json:"language"`
}