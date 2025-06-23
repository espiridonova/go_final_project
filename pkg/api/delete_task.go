package api

import (
	"net/http"

	"github.com/espiridonova/go_final_project/pkg/db"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
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

	err = db.DeleteTask(id)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{err.Error()})
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
