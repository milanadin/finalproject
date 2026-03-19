package api

import (
	"encoding/json"
	"errors"
	"finalproject/pkg/db"
	"log"
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
		writeJson(w, map[string]string{"error": "JSON deserialization error"}, 400)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "task title not specified"}, 400)
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, 400)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "failed to add task"}, 500)
		return
	}

	writeJson(w, map[string]any{"id": strconv.FormatInt(id, 10)}, 200)

}

func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
