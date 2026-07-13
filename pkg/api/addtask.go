package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат JSON"})
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка добавления задачи: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func normalizeDate(task *db.Task) error {
	now := time.Now()
	// Обрезаем время до начала дня для корректного сравнения
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = today.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		if t.Before(today) {
			if task.Repeat == "d 1" {
				task.Date = today.Format(DateFormat)
				return nil
			}
			next, err := NextDate(today, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	} else {
		if t.Before(today) {
			task.Date = today.Format(DateFormat)
		}
	}
	return nil
}