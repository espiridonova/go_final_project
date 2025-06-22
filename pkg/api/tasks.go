package api

import (
	"net/http"

	"github.com/espiridonova/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(search, 50)
	if err != nil {
		writeJson(w, &ErrorResp{"failed to get tasks"})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
