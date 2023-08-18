package main

import (
	"context"
	"fmt"
	"goServer/pkg/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connectToDatabase(database_uri string, l *log.Logger) *pgxpool.Pool {
	// DB
	dbpool, err := pgxpool.New(context.Background(), database_uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	l.Println("Connected to Database")
	return dbpool
}

func main() {
	// logger
	l := log.New(os.Stdout, "Api", log.LstdFlags)

	// env
	err := godotenv.Load(".env")
	if err != nil {
		l.Println(err)
		l.Fatal("Missing environment variables\n")
	}

	// DB
	database_uri := os.Getenv("DATABASE_URL")
	dbpool := connectToDatabase(database_uri, l)
	defer dbpool.Close()

	// Router setup
	sm := chi.NewRouter()
	sm.Use(middleware.Logger)
	sm.Get("/", handlers.HealthCheck)
	sm.Get("/search", func(rw http.ResponseWriter, request *http.Request) {
		handlers.Search(rw, request, dbpool)
	})
	sm.Get("/stock/quote", handlers.GetQuote)
	sm.Get("/stock/candle", handlers.GetCandle)

	// server opts
	server := http.Server{
		Addr:         ":3001",
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// non blocking server
	go func() {
		l.Printf("Listening on port %s\n", server.Addr)

		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	// Notify on Interrupt/kill
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	l.Printf("Received %s, commencing graceful shutdown", <-sigChan)

	// graceful shutdown
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		l.Println(err)
	}
}
