package trade

import (
	"context"
	"time"
)

type ProcessRequest struct {
	FilePaths []string `json:"file_paths" validate:"required,dive,required"`
}

type Processor interface {
	ProcessFiles(files []string) error
}

type Finder interface {
	GetSummary(ctx context.Context, ticker string, fromDate *time.Time) (*Summary, error)
}

type Summary struct {
	Ticker         string  `json:"ticker"`
	MaxRangeValue  float64 `json:"max_range_value"`
	MaxDailyVolume int64   `json:"max_daily_volume"`
}
