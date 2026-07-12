package api

import (
	"net/http"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

// getTaskHandler обрабатывает GET /api/task?id=...
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		WriteJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, task)
}