package api

import (
	"encoding/json"
	"errors"
	"finalproject/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	_, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("invalid repeat format")
		}
	}
	if task.Date < now.Format(DateFormat) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "JSON deserialization error"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "task title not specified"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "failed to add task"})
		return
	}

	writeJson(w, map[string]any{"id": strconv.FormatInt(id, 10)})

}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
