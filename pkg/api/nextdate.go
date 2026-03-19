package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Функция пока поддерживает только "d N" и "y"

	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	arrRepeat := strings.Split(repeat, " ")

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid start date")
	}

	if len(arrRepeat) == 1 && arrRepeat[0] == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	} else if len(arrRepeat) == 2 && arrRepeat[0] == "d" {

		interval, err := strconv.Atoi(arrRepeat[1])

		if err != nil {
			return "", errors.New("repeat incorrect formats")
		}
		if interval < 1 || interval > 365 {
			return "", errors.New("repeat incorrect formats")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil
	} else {
		return "", errors.New("repeat incorrect formats")
	}

}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	var err error
	nowStr := r.FormValue("now")
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "parsing time error", http.StatusBadRequest)
			return
		}
	}
	date, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, "cannot calculate next date", http.StatusBadRequest)
		return
	} else {
		_, err := w.Write([]byte(date))
		if err != nil {
			return
		}
	}
}
