package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"wadood/auth/models"
	"wadood/db"
	"wadood/messages"
	"wadood/utils"
)

type UserService struct{}

func RegisterUser(username, email, password, phoneNumber, imei, language string) (int, *models.User, string, error) {
	if username == "" || email == "" || password == "" || phoneNumber == "" || imei == "" {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("missing_required_fields", language))
	}
	if !utils.EmailRegex.MatchString(email) {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("invalid_email_format", language))
	}
	if !utils.HasComplexity(password) {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("password_too_simple", language))
	}
	if !utils.ValidatePhoneNumber(phoneNumber) {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("invalid_phone_number", language))
	}
	errMsg := ""
	err := utils.ValidatePassword(password)
	if err != nil {
		switch err.Error() {
		case "password must be at least 12 characters long":
			errMsg = messages.GetMessage("password_too_short", language)
		case "password must contain at least one uppercase letter":
			errMsg = messages.GetMessage("password_missing_upper", language)
		case "password must contain at least one number":
			errMsg = messages.GetMessage("password_missing_number", language)
		case "password must contain at least one special character":
			errMsg = messages.GetMessage("password_missing_special", language)
		default:
			errMsg = messages.GetMessage("password_too_simple", language)
		}
		return http.StatusBadRequest, nil, "", fmt.Errorf(errMsg)
	}
	var exists int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?) OR phone_number = ?", email, username, phoneNumber).Scan(&exists)
	if err != nil {
		return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("database_error", language))
	}
	if exists > 0 {
		return http.StatusConflict, nil, "", fmt.Errorf(messages.GetMessage("username_or_email_already_exists", language))
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	res, err := db.DB.Exec("INSERT INTO users (username, email, password, phone_number, tokens, last_login, language) VALUES (?, ?, ?, ?, '{}', ?, ?)", username, email, string(hashed), phoneNumber, time.Now(), language)
	if err != nil {
		return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("registration_failed", language))
	}
	id, _ := res.LastInsertId()
	user := &models.User{ID: int(id), Username: username, Email: email, PhoneNumber: phoneNumber, Tokens: make(map[string]string), LastLogin: time.Now()}
	token, err := utils.GenerateJWT(user.ID, imei)
	if err != nil {
		return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("token_generation_error", language))
	}
	user.Tokens[imei] = token
	tokenJson, _ := json.Marshal(user.Tokens)
	_, err = db.DB.Exec("UPDATE users SET tokens = ? WHERE id = ?", tokenJson, user.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("session_update_error", language))
	}
	return http.StatusCreated, user, token, nil
}

func LoginUser(login, password, imei, language string) (int, *models.User, string, error) {
	if login == "" && password == "" && imei == "" {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("missing_login_password_imei", language))
	}
	if login == "" {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("missing_login", language))
	}
	if password == "" {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("missing_password", language))
	}
	if imei == "" {
		return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("missing_imei", language))
	}
	user, err := getUserByIdentifier(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return http.StatusNotFound, nil, "", fmt.Errorf(messages.GetMessage("user_not_found", language))
		}
		return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("database_error", language))
	}
	if user == nil {
		return http.StatusNotFound, nil, "", fmt.Errorf(messages.GetMessage("user_not_found", language))
	}
	if user.Password == "" {
		return http.StatusUnauthorized, nil, "", fmt.Errorf(messages.GetMessage("user_has_no_password", language))
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return http.StatusUnauthorized, nil, "", fmt.Errorf(messages.GetMessage("invalid_credentials", language))
	}
	token, err := utils.GenerateJWT(user.ID, imei)
	if err != nil {
		return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("token_generation_error", language))
	}
	if user.Tokens == nil {
		user.Tokens = make(map[string]string)
	}
	user.Tokens[imei] = token
	tokenJson, _ := json.Marshal(user.Tokens)
	_, err = db.DB.Exec("UPDATE users SET tokens = ? WHERE id = ?", tokenJson, user.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("session_update_error", language))
	}
	return http.StatusOK, user, token, nil
}

func LogoutUser(userID int, token string, imei string, language string) (int, error) {
	user, err := getUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return http.StatusNotFound, fmt.Errorf(messages.GetMessage("user_not_found", language))
		}
		return http.StatusInternalServerError, fmt.Errorf(messages.GetMessage("db_error", language))
	}

	if user.Tokens == nil || user.Tokens[imei] != token {
		return http.StatusUnauthorized, fmt.Errorf(messages.GetMessage("invalid_token", language))
	}

	delete(user.Tokens, imei)

	tokenJson, err := json.Marshal(user.Tokens)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(messages.GetMessage("token_serialization_error", language))
	}

	_, err = db.DB.Exec("UPDATE users SET tokens = ? WHERE id = ?", tokenJson, user.ID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(messages.GetMessage("logout_failed", language))
	}

	return http.StatusOK, nil
}

func getUserByIdentifier(login string) (*models.User, error) {
	user := &models.User{}
	var tokensStr sql.NullString
	var lastLoginRaw sql.NullString
	query := `SELECT id, username, email, password, tokens, last_login FROM users WHERE LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?) LIMIT 1`
	err := db.DB.QueryRow(query, login, login).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &tokensStr, &lastLoginRaw)
	if err != nil {
		return nil, err
	}
	if lastLoginRaw.Valid && lastLoginRaw.String != "" {
		t, err := time.Parse("2006-01-02 15:04:05", lastLoginRaw.String)
		if err == nil {
			user.LastLogin = t
		}
	}
	if tokensStr.Valid && tokensStr.String != "" {
		err := json.Unmarshal([]byte(tokensStr.String), &user.Tokens)
		if err != nil {
			return nil, err
		}
	} else {
		user.Tokens = make(map[string]string)
	}
	return user, nil
}

func getUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	var tokensStr sql.NullString
	var lastLoginRaw sql.NullString
	query := `SELECT id, username, email, password, tokens, last_login FROM users WHERE id = ? LIMIT 1`

	err := db.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&tokensStr,
		&lastLoginRaw,
	)
	if err != nil {
		return nil, err
	}
	if lastLoginRaw.Valid && lastLoginRaw.String != "" {
		t, err := time.Parse("2006-01-02 15:04:05", lastLoginRaw.String)
		if err == nil {
			user.LastLogin = t
		}
	}
	if tokensStr.Valid && tokensStr.String != "" {
		err := json.Unmarshal([]byte(tokensStr.String), &user.Tokens)
		if err != nil {
			return nil, err
		}
	} else {
		user.Tokens = make(map[string]string)
	}
	return user, nil
}

func FindUserByToken(token string) (*models.User, string, error) {
	rows, err := db.DB.Query("SELECT id, username, email, password, tokens, last_login FROM users")
	if err != nil {
		return nil, "", err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}(rows)
	for rows.Next() {
		var user models.User
		var tokensStr sql.NullString
		var lastLoginRaw sql.NullString
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &tokensStr, &lastLoginRaw)
		if err != nil {
			continue
		}
		if lastLoginRaw.Valid && lastLoginRaw.String != "" {
			t, err := time.Parse("2006-01-02 15:04:05", lastLoginRaw.String)
			if err == nil {
				user.LastLogin = t
			}
		}
		if tokensStr.Valid && tokensStr.String != "" {
			var tokens map[string]string
			err := json.Unmarshal([]byte(tokensStr.String), &tokens)
			if err != nil {
				continue
			}
			for imei, t := range tokens {
				if t == token {
					user.Tokens = tokens
					return &user, imei, nil
				}
			}
		}
	}
	return nil, "", nil
}
