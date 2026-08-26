package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"khhub/internal/auth"
	"khhub/internal/config"
	apphttp "khhub/internal/http"
	"khhub/internal/seed"
	"khhub/internal/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := store.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	q := store.New(pool)
	if err := seedAdmin(ctx, q, cfg); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	var resetSeed func(context.Context) error
	if !cfg.Production() {
		resetSeed = func(c context.Context) error {
			return seed.Reset(c, pool, q)
		}
		n, err := q.CountPublishers(ctx)
		if err != nil {
			log.Fatalf("count publishers: %v", err)
		}
		if n == 0 {
			if err := seed.Reset(ctx, pool, q); err != nil {
				log.Fatalf("demo seed: %v", err)
			}
			log.Printf("loaded demo congregation seed")
		}
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           apphttp.NewRouter(cfg, q, resetSeed),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("khhub api listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func seedAdmin(ctx context.Context, q *store.Queries, cfg config.Config) error {
	n, err := q.CountUsers(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}
	_, err = q.CreateUser(ctx, store.CreateUserParams{
		Email:        cfg.AdminEmail,
		PasswordHash: hash,
	})
	if err != nil {
		return err
	}
	log.Printf("seeded admin user %s — change the password after first login", cfg.AdminEmail)
	return nil
}
