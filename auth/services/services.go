package services

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "wadood/auth/models"
    "wadood/utils"
    "wadood/db"
    "wadood/messages"
    "golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, email, password, imei, language string) (int, *models.User, string, error) {
    if username == "" || email == "" || password == "" || imei == "" {
        return http.StatusBadRequest, nil, "", fmt.Errorf(messages.GetMessage("missing_required_fields", language))
    }
    var exists int
    err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?)", email, username).Scan(&exists)
    if err != nil {
        return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("database_error", language))
    }
    if exists > 0 {
        return http.StatusConflict, nil, "", fmt.Errorf(messages.GetMessage("username_or_email_already_exists", language))
    }
    hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    res, err := db.DB.Exec("INSERT INTO users (username, email, password, tokens, last_login, language) VALUES (?, ?, ?, '{}', ?, ?)", username, email, string(hashed), time.Now(), language)
    if err != nil {
        return http.StatusInternalServerError, nil, "", fmt.Errorf(messages.GetMessage("registration_failed", language))
    }
    id, _ := res.LastInsertId()
    user := &models.User{ID: int(id), Username: username, Email: email, Tokens: make(map[string]string), LastLogin: time.Now()}
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
        if err == sql.ErrNoRows {
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
        if err == sql.ErrNoRows {
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
        json.Unmarshal([]byte(tokensStr.String), &user.Tokens)
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
        json.Unmarshal([]byte(tokensStr.String), &user.Tokens)
    } else {
        user.Tokens = make(map[string]string)
    }
    return user, nil
}