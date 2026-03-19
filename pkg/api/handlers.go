package api

import (
	"encoding/json"
	"finalproject/pkg/db"
	"net/http"
	"time"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, 404)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"}, 200)
		return
	}

	writeJson(w, task, 200)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Задача не найдена"}, 404)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, 400)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "failed to update task"}, 500)
		return
	}

	writeJson(w, map[string]string{}, 200)

}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, 400)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"}, 404)
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "Не удалось удалить задачу"}, 500)
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "cannot calculate next date"}, 400)
			return
		}
		err = db.UpdateDate(next, id)
		if err != nil {
			writeJson(w, map[string]string{"error": "Не удалось обновить дату"}, 500)
			return
		}
	}
	writeJson(w, map[string]string{}, 200)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, 400)
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"}, 404)
		return
	}
	writeJson(w, map[string]string{}, 200)
}
