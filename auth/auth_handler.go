package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"wadood/auth/models"
	"wadood/auth/services"
	"wadood/messages"
	"wadood/utils"
)

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IMEI     string `json:"imei"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: "invalid_request_format"})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	status, user, token, err := services.RegisterUser(req.Username, req.Email, req.Password, req.IMEI, language)
	if err != nil {
		c.JSON(status, models.AuthResponse{Status: status, Message: messages.GetMessage(err.Error(), language)})
		return
	}
	c.JSON(http.StatusCreated, models.AuthResponse{Status: http.StatusCreated, Message: messages.GetMessage("registration_successful", language), User: user, Token: token})
}

func Login(c *gin.Context) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		IMEI     string `json:"imei"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Header("X-Language", "en")
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: "invalid_request_format"})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	status, user, token, err := services.LoginUser(req.Login, req.Password, req.IMEI, language)
	c.Header("X-Language", language)
	if err != nil {
		c.JSON(status, models.AuthResponse{Status: status, Message: messages.GetMessage(err.Error(), language)})
		return
	}
	c.JSON(http.StatusOK, models.AuthResponse{Status: http.StatusOK, Message: messages.GetMessage("login_successful", language), User: user, Token: token})
}

func Logout(c *gin.Context) {
	token := utils.ExtractToken(c.Request)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token[7:]
	}
	var req struct {
		UserID int    `json:"user_id"`
		IMEI   string `json:"imei"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Header("X-Language", "en")
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: "invalid_request_format"})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	status, err := services.LogoutUser(req.UserID, token, req.IMEI, language)
	c.Header("X-Language", language)
	if err != nil {
		c.JSON(status, models.AuthResponse{Status: status, Message: messages.GetMessage(err.Error(), language)})
		return
	}
	c.JSON(http.StatusOK, models.AuthResponse{Status: http.StatusOK, Message: messages.GetMessage("logout_successful", language)})
}
