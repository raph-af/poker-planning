package main

import (
	"fmt"
	"net/http"
	"poker-planning/pkg/forms"
	"poker-planning/pkg/models"
	"strconv"
)

func (app *App) Home(writer http.ResponseWriter, request *http.Request) {
	stories, err := app.Database.GetLatestStories()
	if err != nil {
		app.ServerError(writer, err)
		return
	}

	app.RenderHTML(writer, request, "home.page.html", &HTMLData{Stories: stories})
}

func (app *App) NewStory(writer http.ResponseWriter, request *http.Request) { // form
	app.RenderHTML(writer, request, "new.page.html", &HTMLData{
		Form: &forms.NewStory{},
	})
}

func (app *App) CreateStory(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.ClientError(writer, http.StatusBadRequest)
		return
	}

	form := &forms.NewStory{
		Title:   request.PostForm.Get("title"),
		Content: request.PostForm.Get("content"),
	}

	if !form.IsValid() {
		app.RenderHTML(writer, request, "new.page.html", &HTMLData{Form: form})
		return
	}

	id, err := app.Database.InsertStory(form.Title, form.Content)
	if err != nil {
		app.ServerError(writer, err)
		return
	}

	app.Sessions.Put(request.Context(), "flash", "Your story was saved successfully !")

	http.Redirect(writer, request, fmt.Sprintf("/story/%d", id), http.StatusSeeOther)
}

func (app *App) ShowStory(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.NotFound(writer)
		return
	}

	story, err := app.Database.GetStory(id)
	if err != nil {
		app.ServerError(writer, err)
		return
	}
	if story == nil {
		app.NotFound(writer)
		return
	}

	flash := app.Sessions.PopString(request.Context(), "flash")

	app.RenderHTML(writer, request, "show.page.html", &HTMLData{
		Flash: flash,
		Story: story,
	})
}

func (app *App) SignupUser(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "signup.page.html", &HTMLData{
		Form: &forms.SignupUser{},
	})
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.SignupUser{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !form.IsValid() {
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{
			Form: form,
		})
		return
	}

	err = app.Database.InsertUser(form.Name, form.Email, form.Password)
	if err != nil {
		if err == models.ErrDuplicateEmail {
			form.Failures["Email"] = "Email is already in use"
			app.RenderHTML(w, r, "login.page.html", &HTMLData{Form: form})
			return
		} else {
			app.ServerError(w, err)
		}
	}

	msg := "Your signup was successful. Please log in using your credentials."
	app.Sessions.Put(r.Context(), "flash", msg)
	if err == models.ErrDuplicateEmail {
		form.Failures["Email"] = "Address is already in use"
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{Form: form})
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "login.page.html", &HTMLData{Form: &forms.LoginUser{}})
}

func (app *App) VerifyUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.LoginUser{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !form.IsValid() {
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{
			Form: form,
		})
		return
	}

	currentUserId, err := app.Database.VerifyUser(form.Email, form.Password)
	if err == models.ErrInvalidCredentials {
		form.Failures["Generic"] = "Email or Password is incorrect"
		app.RenderHTML(w, r, "login.page.html", &HTMLData{Form: form})
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}
	app.Sessions.Put(r.Context(), "currentUserID", currentUserId)

	http.Redirect(w, r, "/story/new", http.StatusSeeOther)
}

func (app *App) LogoutUser(w http.ResponseWriter, r *http.Request) {
	app.Sessions.Remove(r.Context(), "currentUserID")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
