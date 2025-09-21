// apps/cp-api/internal/app/wire.go
package app

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"xeed/apps/cp-api/internal/adapter/security"
	"xeed/apps/cp-api/internal/adapter/system"
	"xeed/apps/cp-api/internal/config"
	"xeed/apps/cp-api/internal/http/handlers"
	"xeed/apps/cp-api/internal/repo/pg"
	"xeed/apps/cp-api/internal/routers"
	"xeed/apps/cp-api/internal/usecase"
)

func buildHTTP(ctx context.Context, cfg config.Config) (http.Handler, func(), error) {
	// DB pool
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, func() {}, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, func() {}, err
	}

	cleanup := func() { pool.Close() }

	// repos
	userRepo := pg.NewUserRepositoryPG(pool)

	// adapters
	clock := system.Clock{}
	idgen := system.IDGen{}
	hasher := security.BcryptHasher{}

	// usecases
	userSvc := usecase.NewUserService(userRepo, clock, idgen, hasher)

	// handlers
	userH := handlers.NewUserHandler(userSvc)

	// routers
	handler := routers.InitRouter(userH)
	return handler, cleanup, nil
}
