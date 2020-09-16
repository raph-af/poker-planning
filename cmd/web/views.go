package main

import (
	"bytes"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
	"poker-planning/pkg/models"
	"time"
)

type HTMLData struct {
	Form       interface{}
	Path       string
	Story      *models.Story
	Stories    []*models.Story
	Flash      string
	IsLoggedIn bool
	CSRFToken  string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func (app *App) RenderHTML(w http.ResponseWriter, r *http.Request, page string, data *HTMLData) {
	if data == nil {
		data = &HTMLData{}
	}

	data.Path = r.URL.Path
	data.IsLoggedIn = app.IsLoggedIn(r)

	data.CSRFToken = nosurf.Token(r)

	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
		filepath.Join(app.HTMLDir, page),
	}

	fm := template.FuncMap{
		"humanDate": humanDate,
	}

	ts, err := template.New("").Funcs(fm).ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}
