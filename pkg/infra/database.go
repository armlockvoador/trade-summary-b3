package infra

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"os"
)

func NewPgxPool() (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
}

var Module = fx.Options(
	fx.Provide(
		NewPgxPool,
	),
)
