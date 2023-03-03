package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/Andre-Lopez/snippetbox/internal/models"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	port           string
	staticPath     string
	snippets       *models.SnippetModel
	sessionManager *session.Store
}

func main() {
	// Obtain command line args
	PORT := flag.String("port", ":4000", "HTTP Network Address")
	STATIC_PATH := flag.String("static-path", "./ui/static", "Path of static conent to serve")
	dsn := flag.String("dsn", "web:root@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Set up session manager, will use our mysql DB
	storeDb := mysql.New(mysql.Config{
		Database: "snippetbox",
		Db:       db,
		Table:    "sessions",
		Reset:    true,
	})

	store := session.New(session.Config{
		Storage: storeDb,
	})

	// init our app struct
	application := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		port:           *PORT,
		staticPath:     *STATIC_PATH,
		snippets:       &models.SnippetModel{DB: db},
		sessionManager: store,
	}

	// Init our fiber app
	app := application.routes()

	infoLog.Println("Starting on server", *PORT)
	errorLog.Fatal(app.Listen(*PORT))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
