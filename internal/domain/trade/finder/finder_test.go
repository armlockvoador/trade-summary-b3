package finder_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"negotiation-history-B3/internal/domain/trade"
	"negotiation-history-B3/internal/domain/trade/finder"
	"negotiation-history-B3/mocks"
	"negotiation-history-B3/pkg/repository"
)

func TestFinder_GetSummary(t *testing.T) {
	ctx := context.Background()
	ticker := "WINQ25"
	now := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	t.Run("should return summary successfully", func(t *testing.T) {
		mockRepo := mocks.NewTradeRepository(t)
		f := finder.New(mockRepo)

		mockedReturn := &repository.TradeSummary{
			Ticker:         "WINQ25",
			MaxRangeValue:  140750.00,
			MaxDailyVolume: 7,
		}

		mockRepo.
			On("FindSummary", mock.Anything, ticker, &now).
			Return(mockedReturn, nil).
			Once()

		summary, err := f.GetSummary(ctx, ticker, &now)

		assert.NoError(t, err)
		assert.Equal(t, &trade.Summary{
			Ticker:         mockedReturn.Ticker,
			MaxRangeValue:  mockedReturn.MaxRangeValue,
			MaxDailyVolume: mockedReturn.MaxDailyVolume,
		}, summary)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return sql.ErrNoRows", func(t *testing.T) {
		mockRepo := mocks.NewTradeRepository(t)
		f := finder.New(mockRepo)

		mockRepo.
			On("FindSummary", mock.Anything, ticker, &now).
			Return(nil, sql.ErrNoRows).
			Once()

		summary, err := f.GetSummary(ctx, ticker, &now)

		assert.Nil(t, summary)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("should return generic error", func(t *testing.T) {
		mockRepo := mocks.NewTradeRepository(t)
		f := finder.New(mockRepo)

		mockRepo.
			On("FindSummary", mock.Anything, ticker, &now).
			Return(nil, errors.New("db error")).
			Once()

		summary, err := f.GetSummary(ctx, ticker, &now)

		assert.NotNil(t, summary)
		assert.Error(t, err)
		assert.Empty(t, summary.Ticker)
	})
}
