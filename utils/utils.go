package utils

import (
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "github.com/golang-jwt/jwt/v4"
    "time"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var jwtSecret = []byte("your-very-secure-secret")

func RespondJSON(w http.ResponseWriter, status int, resp interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(resp)
}

func HasComplexity(password string) bool {
    var hasUpper, hasLower, hasNumber, hasSymbol bool
    for _, c := range password {
        switch {
        case 'A' <= c && c <= 'Z':
            hasUpper = true
        case 'a' <= c && c <= 'z':
            hasLower = true
        case '0' <= c && c <= '9':
            hasNumber = true
        case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~", c):
            hasSymbol = true
        }
    }
    return hasUpper && hasLower && hasNumber && hasSymbol
}

func ExtractToken(r *http.Request) string {
    bearer := r.Header.Get("Authorization")
    if strings.HasPrefix(strings.ToLower(bearer), "bearer ") {
        return bearer[7:]
    }
    return ""
}

func ValidatePassword(password string) error {
    if len(password) < 12 {
        return fmt.Errorf("password must be at least 12 characters long")
    }
    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return fmt.Errorf("password must contain at least one uppercase letter")
    }
    if !regexp.MustCompile(`\d`).MatchString(password) {
        return fmt.Errorf("password must contain at least one number")
    }
    if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
        return fmt.Errorf("password must contain at least one special character")
    }
    return nil
}

func GenerateJWT(userID int, imei string) (string, error) {
    claims := jwt.MapClaims{
        "userID": userID,
        "imei": imei,
        "lastLogin": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}