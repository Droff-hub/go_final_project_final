package api

import (
	"net/http"
	"time"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "задача с id "+id+" не найдена" {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{})
		return
	}

	// Вычисляем следующую дату относительно даты задачи
	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат даты"})
		return
	}
	next, err := NextDate(date, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка вычисления следующей даты: " + err.Error()})
		return
	}
	if err := db.UpdateTaskDate(id, next); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}