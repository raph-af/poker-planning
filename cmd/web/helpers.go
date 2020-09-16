package main

import "net/http"

func (app *App) IsLoggedIn(r *http.Request) bool {
    isLoggedIn := app.Sessions.Exists(r.Context(), "currentUserID")

    return isLoggedIn
}
