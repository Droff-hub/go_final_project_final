package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Droff-hub/go_final_project/pkg/api"
	"github.com/golang-jwt/jwt/v5"
)

// Run запускает HTTP-сервер на указанном порту.
func Run() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir := "./web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatal("Папка web не найдена. Убедитесь, что вы запускаете сервер из корня проекта.")
	}

	// Регистрируем API-обработчики
	api.Init()

	// Оборачиваем обработчики middleware для аутентификации
	authHandler := authMiddleware(http.DefaultServeMux)

	// Статика не требует аутентификации
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + port
	log.Printf("Сервер запущен на http://localhost%s\n", addr)
	log.Printf("Отдаём файлы из папки %s", webDir)

	if err := http.ListenAndServe(addr, authHandler); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

// authMiddleware проверяет JWT токен из куки для защищённых API
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем статику и /api/login, /api/signin
		if !strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api/login" || r.URL.Path == "/api/signin" {
			next.ServeHTTP(w, r)
			return
		}

		expectedPass := os.Getenv("TODO_PASSWORD")
		if expectedPass == "" {
			// Аутентификация отключена
			next.ServeHTTP(w, r)
			return
		}

		// Получаем токен из куки
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		// Парсим и проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil // должен совпадать с ключом при генерации
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// еще можно проверить, что хэш в токене соответствует текущему паролю
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		if claims["hash"] != hashPassword(expectedPass) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// hashPassword должна совпадать с той, что в signin.go
func hashPassword(pass string) string {
	h := 0
	for _, c := range pass {
		h = h*31 + int(c)
	}
	return string(rune(h % 100000))
}