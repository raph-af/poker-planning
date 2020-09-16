package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func (app *App) ServerError(writer http.ResponseWriter, err error) {
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
}

func (app *App) ClientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}

func (app *App) NotFound(writer http.ResponseWriter) {
	app.ClientError(writer, http.StatusNotFound)
}
