package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ivcp/polls/internal/data"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
	mutex  sync.Mutex
}

func main() {
	var cfg config
	var app application
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app.logger = logger

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		logger.Fatal(err)
	}
	cfg.port = port
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		logger.Fatal("dsn string not set")
	}
	cfg.db.dsn = dsn
	env := os.Getenv("SERVER_ENV")
	if env == "" {
		logger.Fatal("dsn string not set")
	}
	cfg.env = env

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests persecond")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	app.config = cfg

	db, err := app.connectToDB()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	app.models = data.NewModels(db)

	app.setMetrics(db)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
