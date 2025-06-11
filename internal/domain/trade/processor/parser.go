package processor

import (
	"bufio"
	"context"
	"encoding/csv"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	log "go.uber.org/zap"
	"io"
	"negotiation-history-B3/internal/gen/db"
	envutils "negotiation-history-B3/pkg/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseTXTFile(ctx context.Context, path string, out chan<- db.CreateTradeParams, now time.Time) error {
	file, err := os.Open(path)
	if err != nil {
		log.L().Error("Failed to open file", log.String("path", path), log.Error(err))
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = envutils.GetEnvRune("CSV_DELIMITER", ';')

	if envutils.GetEnvBool("SKIP_HEADER", true) {
		_, _ = reader.Read()
	}
	reader.ReuseRecord = true

	var count int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.L().Warn("Error reading CSV record", log.Error(err))
			continue
		}

		trade, err := mapTXTRecordToTrade(record, now)
		if err != nil {
			log.L().Warn("Error parsing trade record", log.Error(err), log.Strings("record", record))
			continue
		}

		select {
		case out <- trade:
			count++
			if count%10000 == 0 {
				log.L().Info("Processed records", log.Int("count", count), log.String("file", path))
			}
		case <-ctx.Done():
			log.L().Warn("Context cancelled during file parse", log.String("file", path))
			return ctx.Err()
		}
	}

	log.L().Info("Finished reading file", log.String("file", path), log.Int("total_records", count))
	return nil
}

func mapTXTRecordToTrade(record []string, now time.Time) (db.CreateTradeParams, error) {
	priceStr := strings.ReplaceAll(record[3], ",", ".")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return db.CreateTradeParams{}, err
	}

	qty, err := strconv.Atoi(record[4])
	if err != nil {
		return db.CreateTradeParams{}, err
	}

	tradeDate, err := time.Parse("2006-01-02", record[8])
	if err != nil {
		return db.CreateTradeParams{}, err
	}

	var tradePrice pgtype.Numeric
	if err := tradePrice.Scan(strconv.FormatFloat(price, 'f', -1, 64)); err != nil {
		return db.CreateTradeParams{}, err
	}

	return db.CreateTradeParams{
		ID:             uuid.New(),
		CloseTime:      pgtype.Text{String: record[5], Valid: true},
		TradeDate:      pgtype.Timestamptz{Time: tradeDate, Valid: true},
		InstrumentCode: record[1],
		TradePrice:     tradePrice,
		TradeQuantity:  pgtype.Int4{Int32: int32(qty), Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: now, Valid: true},
	}, nil
}
