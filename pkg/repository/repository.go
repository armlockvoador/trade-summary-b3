package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"negotiation-history-B3/internal/gen/db"
	"time"
)

type TradeRepository interface {
	SaveBatchCopyFrom(ctx context.Context, conn *pgx.Conn, trades []db.CreateTradeParams) error
	FindSummary(ctx context.Context, ticker string, fromDate *time.Time) (*TradeSummary, error)
}

type TradeSummary struct {
	Ticker         string  `json:"ticker"`
	MaxRangeValue  float64 `json:"max_range_value"`
	MaxDailyVolume int64   `json:"max_daily_volume"`
}
