package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs"
	_ "github.com/lib/pq"
	"log"
	"poker-planning/pkg/models"
	"time"
)

type Config struct {
	Addr      string
	StaticDir string
	HTMLDir   string
	TlsCert   string
	TlsKey    string
}

func main() {
	config := new(Config)
	flag.StringVar(&config.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&config.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&config.HTMLDir, "html-dir", "./ui/html", "Path to HTML templates")
	flag.StringVar(&config.TlsCert, "tls-cert", "./tls/cert.pem", "Path to TLS certificate")
	flag.StringVar(&config.TlsKey, "tls-key", "./tls/key.pem", "Path to TLS key")
	flag.Parse()

	db := connect()
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = true

	app := &App{
		HTMLDir:   config.HTMLDir,
		StaticDir: config.StaticDir,
		Database:  &models.Database{db},
		Sessions:  sessionManager,
		TlsCert:   config.TlsCert,
		TlsKey:    config.TlsKey,
		Addr:      config.Addr,
	}

	app.RunServer()
}

func connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "postgres")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(50)
	return db
}
