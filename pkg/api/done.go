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
		WriteJSON(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			WriteJSON(w, map[string]string{"error": err.Error()})
			return
		}
		WriteJSON(w, map[string]string{})
		return
	}

	// Если правило повторения есть, вычисляем следующую дату относительно даты задачи
	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		WriteJSON(w, map[string]string{"error": "Неверный формат даты"})
		return
	}
	next, err := NextDate(date, task.Date, task.Repeat)
	if err != nil {
		WriteJSON(w, map[string]string{"error": "Ошибка вычисления следующей даты: " + err.Error()})
		return
	}
	if err := db.UpdateTaskDate(id, next); err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, map[string]string{})
}