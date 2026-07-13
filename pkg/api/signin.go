package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// signinHandler обрабатывает POST /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
		return
	}

	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Пароль не задан"})
		return
	}

	if req.Password != expected {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	// Генерируем JWT токен с полезной нагрузкой (хэш пароля)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hashPassword(expected),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка создания токена"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}


// hashPassword возвращает простую контрольную сумму пароля
func hashPassword(pass string) string {
	h := 0
	for _, c := range pass {
		h = h*31 + int(c)
	}
	return string(rune(h % 100000))
}