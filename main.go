package main

import (
	"context"
	"fmt"
	"goServer/handlers"
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

func connectToDatabase(database_uri string) *pgxpool.Pool {
	// DB
	dbpool, err := pgxpool.New(context.Background(), database_uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to Database")
	return dbpool
}

func main() {
	// logger
	l := log.New(os.Stdout, "Api", log.LstdFlags)

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Missing environment variables")
		os.Exit(1)
		return
	}
	database_uri := os.Getenv("DATABASE_URL")
	dbpool := connectToDatabase(database_uri)
	defer dbpool.Close()

	// Router setup
	sm := chi.NewRouter()
	sm.Use(middleware.Logger)
	sm.Get("/search", func(rw http.ResponseWriter, request *http.Request) {
		handlers.Search(rw, request, dbpool)
	})
	sm.Get("/stock/quote", handlers.GetQuote)

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
		fmt.Printf("Listening on port %s\n", server.Addr)

		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	// Notify on Interrupt/kill
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Println("Graceful shutdown", sig)

	// graceful shutdown
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(tc)
}
