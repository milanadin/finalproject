package main

import (
	"finalproject/pkg/api"
	"finalproject/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		fmt.Println("database initialization error", err)
		return
	}
	defer db.Close()

	api.Init()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
