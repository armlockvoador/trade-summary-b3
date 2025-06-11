package finder

import (
	"context"
	"negotiation-history-B3/internal/domain/trade"
	"negotiation-history-B3/pkg/repository"
	"time"
)

type Finder struct {
	repo repository.TradeRepository
}

func New(repo repository.TradeRepository) *Finder {
	return &Finder{
		repo: repo,
	}
}

func (f Finder) GetSummary(ctx context.Context, ticker string, fromDate *time.Time) (*trade.Summary, error) {
	summary, err := f.repo.FindSummary(ctx, ticker, fromDate)
	if err != nil {
		return nil, err

	}
	return &trade.Summary{
		Ticker:         summary.Ticker,
		MaxRangeValue:  summary.MaxRangeValue,
		MaxDailyVolume: summary.MaxDailyVolume,
	}, nil
}
