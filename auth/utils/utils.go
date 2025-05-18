package utils

import (
    "encoding/json"
    "net/http"
    "regexp"
    "strings"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
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
