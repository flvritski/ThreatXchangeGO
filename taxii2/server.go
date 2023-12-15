package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"opentaxii/internal/driver"
	"opentaxii/internal/models"

	"gopkg.in/yaml.v2"
)

const version = "1.0.0"

type TAXII2 struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Contact     string    `yaml:"contact"`
	Default     string    `yaml:"default"`
	ApiRoots    []APIRoot `yaml:"api_roots"`
}

type APIRoot struct {
	Title            string `yaml:"title"`
	Description      string `yaml:"description"`
	URL              string `yaml:"url"`
	MaxContentLength int    `yaml:"max_content_length"`
}

type config struct {
	TAXII TAXII2 `yaml:"taxii2"`
	port  int
	env   string
	db    struct {
		dsn string
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	app.infoLog.Printf("Starting Back end server in %s mode on port %d\n", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	data, err := os.ReadFile("../config/taxii2.yaml")
	if err != nil {
		panic(err)
	}

	// Unmarshal YAML into Config struct
	var cfg config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)

	}

	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application enviroment {development|production|mainteinance}")
	flag.StringVar(&cfg.db.dsn, "dsn", "godeus:secret@tcp(localhost:3306)/godtis?parseTime=true&tls=false", "DSN")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenSQL(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: conn},
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
