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

// SignUserID создает подпись для идентификатора пользователя.
func SignUserID(userID string) string {
	h := hmac.New(sha256.New, []byte(SecretKey))
	h.Write([]byte(userID))
	return hex.EncodeToString(h.Sum(nil))
}

// SetUserCookie устанавливает пользователю подписанную куку с его идентификатором.
func SetUserCookie(c echo.Context, userID string) {
	signature := SignUserID(userID)
	cookieValue := userID + "." + signature
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    cookieValue,
		Expires:  time.Now().Add(24 * time.Hour * 365), // Кука действует 1 год
		HttpOnly: true,
		Secure:   false,
	}
	c.SetCookie(cookie)
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

	userID := parts[0]
	signature := parts[1]

	expectedSignature := SignUserID(userID)
	if hmac.Equal([]byte(expectedSignature), []byte(signature)) {
		return userID, true
	}

	return "", false
}
