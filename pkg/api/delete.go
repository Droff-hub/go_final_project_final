package api

import (
	"net/http"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

// deleteTaskHandler обрабатывает DELETE /api/task?id=...
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		WriteJSON(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, map[string]string{})
}