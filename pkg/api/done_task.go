package api

import (
	"net/http"
	"time"

	"github.com/espiridonova/go_final_project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, &ErrorResp{"Не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{err.Error()})
		return
	}
	if task == nil {
		writeJson(w, http.StatusNotFound, &ErrorResp{"Задача не найдена"})
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, &ErrorResp{err.Error()})
			return
		}
		writeJson(w, http.StatusOK, struct{}{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{err.Error()})
		return
	}

	err = db.UpdateDate(next, id)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{err.Error()})
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
