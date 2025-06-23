package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/espiridonova/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	id, err := db.AddTask(task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, &ErrorResp{err.Error()})
		return
	}

	writeJson(w, http.StatusBadRequest, &ShortTask{id})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}
	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return err
	}
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}
	return nil
}
