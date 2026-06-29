package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
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

	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	log.Println("Running database migrations...")
	runDBMigrations(db)
	log.Println("Database migrations complete")

	catalogRepo := catalogAdapters.NewPostgresRepository(db)
	catalogService := catalog.NewService(catalogRepo)
	catalogHandler := catalogAdapters.NewHTTPHandler(catalogService)

	ordersRepo := ordersAdapters.NewPostgresRepository(db)
	ordersService := orders.NewService(ordersRepo, catalogService)
	ordersHandler := ordersAdapters.NewHTTPHandler(ordersService, cfg.AdminUser, cfg.AdminPass)

	mux := http.NewServeMux()
	catalogHandler.RegisterRoutes(mux)
	ordersHandler.RegisterRoutes(mux)

	// Required for Kubernetes Liveness/Readiness Probes
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Failed to write health response: %v", err)
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	// Give active connections 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func runDBMigrations(db *sql.DB) {
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("Failed to create migration source driver: %v", err)
	}

	dbDriver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		log.Fatalf("Failed to create database driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}
