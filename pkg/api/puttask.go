package api

import (
	"encoding/json"
	"net/http"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат JSON"})
		return
	}

	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}
	if err := normalizeDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "задача с id "+task.ID+" не найдена" {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}