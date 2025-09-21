package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"xeed/apps/cp-api/internal/domain"
	"xeed/apps/cp-api/internal/http/handlers"
	"xeed/apps/cp-api/internal/repo/pg"
	"xeed/apps/cp-api/internal/routers"
	"xeed/apps/cp-api/internal/usecase"
)

type sysClock struct{}

func (sysClock) Now() time.Time { return time.Now().UTC() }

type sysIDGen struct{}

func (sysIDGen) New() uuid.UUID { return uuid.New() }

// NOTE: hanya untuk test. Ganti dengan bcrypt di produksi.
type fakeHasher struct{}

func (fakeHasher) Hash(plain string) (string, domain.PasswordAlg, time.Time, error) {
	return plain, domain.PasswordAlg("plain"), time.Now().UTC(), nil
}

func Run() error {
	// ENV
	_ = godotenv.Load(".env")

	port := getenv("PORT", "8080")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is required")
	}

	// DB pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	if err := pool.Ping(ctx); err != nil {
		return err
	}
	defer pool.Close()

	// ===== Wiring (DI) =====
	userRepo := pg.NewUserRepositoryPG(pool)
	// isi adapter nyata sesuai kontrakmu (clock, idgen, hasher) di service.NewUserService
	// >>> JANGAN pakai 'var clock contract.Clock' lalu 'clock :=' di bawah
	clock := sysClock{} // <- perhatikan: TIDAK ada deklarasi 'var clock ...' di atas
	idgen := sysIDGen{}
	hasher := fakeHasher{}
	userSvc := usecase.NewUserService(
		userRepo,
		clock,  // contract.Clock
		idgen,  // contract.IDGen
		hasher, // contract.PasswordHasher
	)
	userHandler := handlers.NewUserHandler(userSvc)

	// ===== Router =====
	r := routers.InitRouter(userHandler)

	// ===== HTTP Server =====
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	// Start
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	select {
	case <-quit:
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctxShutdown)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
