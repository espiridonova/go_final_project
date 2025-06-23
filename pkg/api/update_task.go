package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/espiridonova/go_final_project/pkg/db"
)

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{"internal server error"})
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{"error unmarshal"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, &ErrorResp{"title is empty"})
		return
	}

	err = checkDate(task)
	if err != nil {
		msg := fmt.Sprintf("invalid task date: %s; err: %s", task.Date, err)
		writeJson(w, http.StatusBadRequest, &ErrorResp{msg})
		return
	}

	err = db.UpdateTask(task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{err.Error()})
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
