package api

import (
	"encoding/json"
	"net/http"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

// putTaskHandler обрабатывает PUT /api/task
func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		WriteJSON(w, map[string]string{"error": "Неверный формат JSON"})
		return
	}

	if task.ID == "" {
		WriteJSON(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}
	if task.Title == "" {
		WriteJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}
	// Нормализация даты 
	if err := normalizeDate(&task); err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, map[string]string{}) // пустой JSON
}