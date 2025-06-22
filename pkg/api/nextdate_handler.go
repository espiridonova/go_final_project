package api

import (
	"net/http"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(dateLayout, nowParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
