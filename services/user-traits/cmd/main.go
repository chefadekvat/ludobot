package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"user-traits/gen/api"
	"user-traits/internal/di"
	"user-traits/internal/server"

	"github.com/akamensky/argparse"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Arguments struct {
	psqlConnString string
	port           string
}

func parseArgs() Arguments {
	parser := argparse.NewParser("user-traits", "User traits service for managing user data")

	psqlConnString := parser.String("", "postgresql", &argparse.Options{
		Required: true,
		Help:     "PSQL connection string",
	})
	port := parser.String("", "port", &argparse.Options{
		Required: true,
		Help:     "Server port",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return Arguments{
		psqlConnString: *psqlConnString,
		port:           *port,
	}
}

func createDBPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func run() error {
	arguments := parseArgs()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pool, err := createDBPool(ctx, arguments.psqlConnString)
	if err != nil {
		return fmt.Errorf("database pool creation error: %w", err)
	}
	defer pool.Close()

	depsFactory := di.NewDependenciesFactory(pool)
	srv := server.NewServer(depsFactory)

	handler := api.NewStrictHandler(srv, nil)
	httpHandler := api.Handler(handler)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", arguments.port),
		Handler: httpHandler,
	}

	slog.Info("Starting HTTP server", "port", arguments.port)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("HTTP server error: %s", err.Error()))
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down server...")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
