package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "go.uber.org/zap"
	"negotiation-history-B3/internal/gen/db"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewTradeRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r Repository) SaveBatchCopyFrom(ctx context.Context, conn *pgx.Conn, batch []db.CreateTradeParams) error {
	_, err := conn.CopyFrom(ctx,
		pgx.Identifier{"trade"},
		[]string{"id", "close_time", "trade_date", "instrument_code", "trade_price", "trade_quantity", "created_at", "updated_at", "deleted"},
		pgx.CopyFromSlice(len(batch), func(i int) ([]interface{}, error) {
			t := batch[i]
			return []interface{}{
				t.ID,
				t.CloseTime,
				t.TradeDate,
				t.InstrumentCode,
				t.TradePrice,
				t.TradeQuantity,
				t.CreatedAt,
				t.UpdatedAt,
				false,
			}, nil
		}),
	)
	return err
}

func (r Repository) FindSummary(ctx context.Context, ticker string, fromDate *time.Time) (*TradeSummary, error) {
	baseQuery := `
		SELECT $1 AS ticker, max_price.max_price, max_volume_day.max_volume
		FROM (
			SELECT MAX(CAST(trade_price AS NUMERIC)) AS max_price
			FROM trade
			WHERE instrument_code = $1
			%s
		) AS max_price,
		(
			SELECT MAX(volume) AS max_volume
			FROM (
				SELECT trade_date, SUM(trade_quantity) AS volume
				FROM trade
				WHERE instrument_code = $1
				%s
				GROUP BY trade_date
			) daily_volume
		) max_volume_day
	`

	var (
		query string
		args  []any
	)

	if fromDate != nil {
		dateFilter := "AND trade_date >= $2"
		query = fmt.Sprintf(baseQuery, dateFilter, dateFilter)
		args = []any{ticker, *fromDate}
	} else {
		query = fmt.Sprintf(baseQuery, "", "")
		args = []any{ticker}
	}

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.L().Error("Failed to acquire DB connection", log.Error(err))
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, query, args...)

	var summary TradeSummary
	if err := row.Scan(&summary.Ticker, &summary.MaxRangeValue, &summary.MaxDailyVolume); err != nil {
		if err.Error() == "no rows in result set" {
			log.L().Warn("No data found for ticker", log.String("ticker", ticker))
			return &summary, nil
		}
		log.L().Error("Failed to scan summary result", log.Error(err), log.String("ticker", ticker))
		return nil, err
	}

	log.L().Info("Successfully retrieved trade summary", log.String("ticker", ticker))
	return &summary, nil
}
