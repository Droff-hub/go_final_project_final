package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

// addTaskHandler обрабатывает POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		WriteJSON(w, map[string]string{"error": "Неверный формат JSON"})
		return
	}

	if task.Title == "" {
		WriteJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := normalizeDate(&task); err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		WriteJSON(w, map[string]string{"error": "Ошибка добавления задачи: " + err.Error()})
		return
	}

	WriteJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func normalizeDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}
	// Если правило не пустое, проверяем его корректность через NextDate (от сегодня)
	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err // неверное правило
		}
	}
	// Если дата в прошлом, ставим сегодня (независимо от правила)
	if t.Before(now) {
		task.Date = now.Format(DateFormat)
	}
	return nil
}