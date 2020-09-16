package main

import (
	"github.com/alexedwards/scs"
	"poker-planning/pkg/models"
)

type App struct {
	HTMLDir   string
	StaticDir string
	Database  *models.Database
	Sessions  *scs.SessionManager
	TlsCert   string
	TlsKey    string
	Addr      string
}
