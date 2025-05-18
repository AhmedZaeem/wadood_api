package auth

import (
	"encoding/json"
	"net/http"

	"wadood/auth/models"
	"wadood/auth/services"
	"wadood/auth/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, models.AuthResponse{Message: services.GetMessage("invalid_request", req.Language)})
		return
	}
	resp, status := services.RegisterUser(&req)
	utils.RespondJSON(w, status, resp)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, models.AuthResponse{Message: services.GetMessage("invalid_request", req.Language)})
		return
	}
	resp, status := services.LoginUser(&req)
	utils.RespondJSON(w, status, resp)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := utils.ExtractToken(r)
	if token == "" {
		utils.RespondJSON(w, http.StatusUnauthorized, models.AuthResponse{Message: services.GetMessage("unauthorized", "en")})
		return
	}
	resp, status := services.LogoutUser(token)
	utils.RespondJSON(w, status, resp)
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	token := utils.ExtractToken(r)
	if token == "" {
		utils.RespondJSON(w, http.StatusUnauthorized, models.AuthResponse{Message: services.GetMessage("unauthorized", "en")})
		return
	}
	var req models.EditProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, models.AuthResponse{Message: services.GetMessage("invalid_request", req.Language)})
		return
	}
	resp, status := services.EditUserProfile(token, &req)
	utils.RespondJSON(w, status, resp)
}