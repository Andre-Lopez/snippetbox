package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net"
	"os"

	"github.com/Andre-Lopez/snippetbox/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	port           string
	staticPath     string
	snippets       models.SnippetModelInterface
	sessionManager *session.Store
	users          models.UserModelInterface
}

func main() {
	mux, app, ln, db, PORT, err := Setup()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app.infoLog.Println("Starting on server", *PORT)
	app.errorLog.Fatal(mux.Listener(ln))
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

func Setup() (*fiber.App, *application, net.Listener, *sql.DB, *string, error) {
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
		return nil, nil, nil, nil, nil, err
	}

	// Set up session manager, will use our mysql DB
	storeDb := mysql.New(mysql.Config{
		Database: "snippetbox",
		Db:       db,
		Table:    "sessions",
		Reset:    false,
	})

	store := session.New(session.Config{
		Storage:      storeDb,
		CookieSecure: true,
	})

	// init our app struct
	application := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		port:           *PORT,
		staticPath:     *STATIC_PATH,
		snippets:       &models.SnippetModel{DB: db},
		sessionManager: store,
		users:          &models.UserModel{DB: db},
	}

	// Init our fiber app
	mux := application.routes()

	// Set up TLS Cert
	cer, err := tls.LoadX509KeyPair("tls/cert.pem", "tls/key.pem")
	if err != nil {
		errorLog.Fatal(err)
		return nil, nil, nil, nil, nil, err
	}

	config := &tls.Config{
		Certificates:     []tls.Certificate{cer},
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	ln, err := tls.Listen("tcp", *PORT, config)
	if err != nil {
		panic(err)
	}

	return mux, application, ln, db, PORT, nil
}
