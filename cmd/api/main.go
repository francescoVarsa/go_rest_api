package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config = struct {
	port int
	env  string
}

type AppStatus struct {
	Environment string `json:"environment"`
	Status      string `json:"status"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "developement", "Application environment (developement | production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	fmt.Println("Running")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second}

	logger.Println("Starting server on port", cfg.port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
