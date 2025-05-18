package services

import (
    "encoding/json"
    "time"
    "wadood/messages"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "fmt"
    "wadood/auth/utils"
    "wadood/auth/models"
    "wadood/auth/db"
)

var SecretKey = "your_secret_key"

func RegisterUser(reqUser models.User, password, imei string) (*models.User, string, error) {
    if !utils.EmailRegex.MatchString(reqUser.Email) {
        return nil, "", fmt.Errorf(messages.Get("Invalid email format", reqUser.Language))
    }
    if len(password) < 12 || !utils.HasComplexity(password) {
        return nil, "", fmt.Errorf(messages.Get("Password must be at least 12 characters and include upper, lower, number, and symbol.", reqUser.Language))
    }
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, "", err
    }
    token := GenerateToken(reqUser.Email, imei)
    tokens := map[string]string{imei: token}
    jsonTokens, _ := json.Marshal(tokens)
    res, err := db.DB.Exec(`INSERT INTO users (username, email, password, tokens, last_login, language) VALUES (?, ?, ?, ?, ?, ?)`,
        reqUser.Username, reqUser.Email, string(hashed), string(jsonTokens), time.Now(), reqUser.Language)
    if err != nil {
        return nil, "", err
    }
    id, _ := res.LastInsertId()
    user := &models.User{
        ID:        int(id),
        Username:  reqUser.Username,
        Email:     reqUser.Email,
        Tokens:    tokens,
        LastLogin: time.Now(),
        Language:  reqUser.Language,
    }
    return user, token, nil
}

func GenerateToken(email, imei string) string {
    claims := jwt.MapClaims{
        "email": email,
        "imei":  imei,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(SecretKey))
    return tokenString
}