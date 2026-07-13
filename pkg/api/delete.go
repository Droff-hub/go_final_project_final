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
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "задача с id "+id+" не найдена" {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}