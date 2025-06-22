package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
	http.HandleFunc("/api/signin", signInHandler)

}

type ErrorResp struct {
	Error string `json:"error"`
}

type ShortTask struct {
	ID int64 `json:"id"`
}

func writeJson(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("error:", err)
	}
}
