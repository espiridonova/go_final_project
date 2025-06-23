package api

import (
	"net/http"

	"github.com/espiridonova/go_final_project/pkg/db"
)

const (
	maxLimit = 50
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(search, maxLimit)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{"failed to get tasks"})
		return
	}
	writeJson(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
