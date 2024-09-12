package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	CookieName = "user_id"
	SecretKey  = "super-secret-key" // Это ключ для подписи куки
)

// SignUser создает подпись для идентификатора пользователя.
func SignUser(userLogin string) string {
	h := hmac.New(sha256.New, []byte(SecretKey))
	h.Write([]byte(userLogin))
	return hex.EncodeToString(h.Sum(nil))
}

// GetUserCookie устанавливает пользователю подписанную куку с его идентификатором.
func GetUserCookie(userLogin string) *http.Cookie {
	signature := SignUser(userLogin)
	cookieValue := userLogin + "." + signature
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    cookieValue,
		Expires:  time.Now().Add(24 * time.Hour * 365), // Кука действует 1 год
		HttpOnly: true,
		Secure:   false,
	}
	return cookie
}

// ValidateUserCookie проверяет подлинность куки и возвращает идентификатор пользователя.
func ValidateUserCookie(c echo.Context) (string, bool) {
	cookie, err := c.Cookie(CookieName)
	if err != nil {
		return "", false
	}

	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return "", false
	}

	userLogin := parts[0]
	signature := parts[1]

	expectedSignature := SignUser(userLogin)
	if hmac.Equal([]byte(expectedSignature), []byte(signature)) {
		return userLogin, true
	}

	return "", false
}
