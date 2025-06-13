package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"wadood/auth/models"
	"wadood/auth/services"
	"wadood/messages"
	utils "wadood/utils"
)

func Register(c *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		IMEI        string `json:"imei"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("invalid_json_format", "en")})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	if req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("missing_phone_number", language)})
		return
	}
	if !utils.ValidatePhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("invalid_phone_number", language)})
		return
	}
	regData := map[string]string{
		"username":     req.Username,
		"email":        req.Email,
		"password":     req.Password,
		"phone_number": req.PhoneNumber,
		"imei":         req.IMEI,
		"language":     language,
	}
	SetRegistrationData(req.PhoneNumber, regData)
	otp := generateOTP()
	SetOTP(req.PhoneNumber, otp)
	err := utils.SendOTP(req.PhoneNumber, fmt.Sprintf("Your OTP code is: %s", otp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.AuthResponse{Status: http.StatusInternalServerError, Message: messages.GetMessage("otp_failed", language)})
		return
	}
	c.JSON(http.StatusOK, models.AuthResponse{Status: http.StatusOK, Message: messages.GetMessage("otp_sent", language)})
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

func SendOTP(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_json_format", "en")})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	if !utils.ValidatePhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_phone_number", language)})
		return
	}
	otp := generateOTP()
	SetOTP(req.PhoneNumber, otp)
	fmt.Printf("OTP for %s is: %s\n", req.PhoneNumber, otp)
	c.JSON(http.StatusOK, gin.H{"message": messages.GetMessage("otp_sent", language)})
}

func VerifyOTP(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number"`
		OTP         string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_json_format", "en")})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	storedOTP, ok := GetOTP(req.PhoneNumber)
	if !ok || storedOTP != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"message": messages.GetMessage("otp_invalid", language)})
		return
	}
	regData, found := GetRegistrationData(req.PhoneNumber)
	if !found {
		DeleteOTP(req.PhoneNumber)
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("missing_required_fields", language)})
		return
	}
	status, user, token, err := services.RegisterUser(
		regData["username"],
		regData["email"],
		regData["password"],
		regData["phone_number"],
		regData["imei"],
		regData["language"],
	)
	DeleteOTP(req.PhoneNumber)
	DeleteRegistrationData(req.PhoneNumber)
	if err != nil {
		c.JSON(status, gin.H{"message": messages.GetMessage(err.Error(), language)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": messages.GetMessage("registration_successful", language), "user": user, "token": token})
}

func ResendOTP(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("invalid_json_format", "en")})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	_, found := GetRegistrationData(req.PhoneNumber)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"message": messages.GetMessage("missing_required_fields", language)})
		return
	}
	otp := generateOTP()
	SetOTP(req.PhoneNumber, otp)
	err := utils.SendOTP(req.PhoneNumber, fmt.Sprintf("Your OTP code is: %s", otp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": messages.GetMessage("otp_failed", language)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": messages.GetMessage("otp_sent", language)})
}

func ForgetPassword(c *gin.Context) {
	var req struct {
		Login            string `json:"login"`
		VerificationType string `json:"verification_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("invalid_json_format", "en")})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	var user *models.User
	var err error
	if req.Login != "" {
		user, err = services.FindUserByLogin(req.Login)
	} else {
		err = fmt.Errorf("login field is empty")
	}
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, models.AuthResponse{Status: http.StatusNotFound, Message: messages.GetMessage("forget_password_user_not_found", language)})
		return
	}
	otp := generateOTP()
	SetOTP(req.Login, otp)
	var sendErr error
	switch req.VerificationType {
	case "phone", "whatsapp":
		// Always validate phone number before sending SMS
		if !utils.ValidatePhoneNumber(user.PhoneNumber) {
			c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("invalid_phone_number", language)})
			return
		}
		sendErr = utils.SendOTP(user.PhoneNumber, fmt.Sprintf("Your OTP code is: %s", otp))
	case "email":
		if !utils.EmailRegex.MatchString(user.Email) {
			c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("invalid_email_format", language)})
			return
		}
		//subject := "Your Password Reset OTP"
		//sendErr = utils.SendMailgunEmail(user.Email, subject, otp)
	default:
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("invalid_input", language)})
		return
	}
	if sendErr != nil {
		c.JSON(http.StatusInternalServerError, models.AuthResponse{Status: http.StatusInternalServerError, Message: messages.GetMessage("forget_password_failed", language)})
		return
	}
	c.JSON(http.StatusOK, models.AuthResponse{Status: http.StatusOK, Message: messages.GetMessage("forget_password_otp_sent", language)})
}

func VerifyForgetPasswordOTP(c *gin.Context) {
	var req struct {
		Login string `json:"login"`
		OTP   string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{Status: http.StatusBadRequest, Message: messages.GetMessage("invalid_json_format", "en")})
		return
	}
	language := c.GetHeader("X-Language")
	if language == "" {
		language = "en"
	}
	storedOTP, ok := GetOTP(req.Login)
	if !ok || storedOTP != req.OTP {
		c.JSON(http.StatusUnauthorized, models.AuthResponse{Status: http.StatusUnauthorized, Message: messages.GetMessage("otp_invalid", language)})
		return
	}
	DeleteOTP(req.Login)
	c.JSON(http.StatusOK, models.AuthResponse{Status: http.StatusOK, Message: messages.GetMessage("otp_verified", language)})
}
