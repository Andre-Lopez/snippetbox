package main

import (
	"flag"
	"log"
	"os"
)

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	port       string
	staticPath string
}

func main() {
	// Obtain command line args
	PORT := flag.String("port", ":4000", "HTTP Network Address")
	STATIC_PATH := flag.String("static-path", "./ui/static", "Path of static conent to serve")
	flag.Parse()

	// Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// init our app struct
	application := &application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		port:       *PORT,
		staticPath: *STATIC_PATH,
	}

	// Init our fiber app
	app := application.routes()

	infoLog.Println("Starting on server", *PORT)
	errorLog.Fatal(app.Listen(*PORT))
}
