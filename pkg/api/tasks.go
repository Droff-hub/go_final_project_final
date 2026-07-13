package api

import (
	"net/http"

	"github.com/Droff-hub/go_final_project/pkg/db"
)

// TasksResponse структура ответа со списком задач.
type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр search из запроса (это для задания со звёздочкой)
	search := r.FormValue("search")

	// Лимит можно сделать константой (по стандарту 50)
	const limit = 50

	tasks, err := db.Tasks(limit, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, TasksResponse{Tasks: tasks})
}