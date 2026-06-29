package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmorall/booksonline/db/migrations"
	"github.com/mmorall/booksonline/internal/catalog"
	catalogAdapters "github.com/mmorall/booksonline/internal/catalog/adapters"
	"github.com/mmorall/booksonline/internal/config"
	"github.com/mmorall/booksonline/internal/orders"
	ordersAdapters "github.com/mmorall/booksonline/internal/orders/adapters"

	_ "github.com/mmorall/booksonline/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

// @title BooksOnline API
// @version 1.0
// @description Backend service for managing catalog and orders.
// @host api-booksonline.miguelmoral.com
// @BasePath /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error closing database connection", "error", err)
		}
	}()

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to PostgreSQL")

	slog.Info("Running database migrations...")
	runDBMigrations(db)
	slog.Info("Database migrations complete")

	catalogRepo := catalogAdapters.NewPostgresRepository(db)
	catalogService := catalog.NewService(catalogRepo)
	catalogHandler := catalogAdapters.NewHTTPHandler(catalogService)

	ordersRepo := ordersAdapters.NewPostgresRepository(db)
	ordersService := orders.NewService(ordersRepo, catalogService)
	ordersHandler := ordersAdapters.NewHTTPHandler(ordersService, cfg.AdminUser, cfg.AdminPass)

	mux := http.NewServeMux()

	catalogHandler.RegisterRoutes(mux)
	ordersHandler.RegisterRoutes(mux)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// Required for Kubernetes Liveness/Readiness Probes
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			slog.Error("Failed to write health response", "error", err)
		}
	})

	handlerWithCORS := corsMiddleware(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlerWithCORS,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		slog.Info("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "error", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server gracefully...")

	// Give active connections 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited properly")
}

func runDBMigrations(db *sql.DB) {
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		slog.Error("Failed to create migration source driver", "error", err)
		os.Exit(1)
	}

	dbDriver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		slog.Error("Failed to create database driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		slog.Error("Failed to initialize migrator", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://booksonline.miguelmoral.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
