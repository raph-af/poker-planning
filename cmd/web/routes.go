package main

import (
    "github.com/bmizerany/pat"
    "net/http"
)

func (app *App) Routes() http.Handler {
    mux := pat.New()
    mux.Get("/", http.HandlerFunc(app.Home))

    mux.Get("/story/new", http.HandlerFunc(app.NewStory))
    mux.Post("/story/new", http.HandlerFunc(app.CreateStory))
    mux.Get("/story/:id", NoSurf(http.HandlerFunc(app.ShowStory)))

    mux.Get("/user/signup", NoSurf(http.HandlerFunc(app.SignupUser)))
    mux.Post("/user/signup", NoSurf(http.HandlerFunc(app.CreateUser)))
    mux.Get("/user/login", NoSurf(http.HandlerFunc(app.LoginUser)))
    mux.Post("/user/login", NoSurf(http.HandlerFunc(app.VerifyUser)))
    mux.Post("/user/logout", app.RequireLogin(http.HandlerFunc(app.LogoutUser)))

    fileServer := http.FileServer(http.Dir(app.StaticDir))
    mux.Get("/static/", http.StripPrefix("/static", fileServer))

    return LogRequest(SecureHeaders(app.Sessions.LoadAndSave(mux)))
}
