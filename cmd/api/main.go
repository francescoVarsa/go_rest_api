package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config = struct {
	port int
	env  string
	db   struct {
		dsn string
	}
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
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://francesco:{PASSWORD}@localhost/go_movies?sslmode=disable", "Postgres connection string")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pwd, err := getPwdFromFile("./secrets/passwords.txt")

	if err != nil {
		logger.Fatal(err, "Cannot get password")
	}

	db, err := openDB(cfg, pwd)

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

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

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config, password string) (*sql.DB, error) {
	db, err := sql.Open("postgres", strings.Replace(cfg.db.dsn, "{PASSWORD}", password, 1))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil

}

func getPwdFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	pwd := ""

	for scanner.Scan() {
		ln := scanner.Text()

		splittedStr := strings.Split(ln, ":")

		if splittedStr[0] == "database_password" {
			pwd = strings.Trim(splittedStr[1], " ")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return pwd, nil
}
