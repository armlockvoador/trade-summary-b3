package processor

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	log "go.uber.org/zap"
	"negotiation-history-B3/internal/gen/db"
	"negotiation-history-B3/pkg/repository"
	envutils "negotiation-history-B3/pkg/utils"
	"runtime"
	"sync"
	"time"
)

type Processor struct {
	repo repository.TradeRepository
	pool *pgxpool.Pool
}

func New(repo repository.TradeRepository, pool *pgxpool.Pool) *Processor {
	return &Processor{
		repo: repo,
		pool: pool,
	}
}

func (p Processor) ProcessFiles(files []string) error {
	ctx := context.Background()

	batchSize := envutils.GetEnvInt("BATCH_SIZE", 1000)
	numWorkers := envutils.GetEnvInt("NUM_WORKERS", runtime.NumCPU()*2)
	bufferSize := envutils.GetEnvInt("MAX_CHANNEL_BUFFER", 10000)
	tickerSeconds := envutils.GetEnvInt("TICKER_SECONDS", 2)

	tradeChan := make(chan db.CreateTradeParams, bufferSize)
	errChan := make(chan error, numWorkers)

	log.L().Info("Starting workers",
		log.Int("num_workers", numWorkers),
		log.Int("batch_size", batchSize),
		log.Int("buffer_size", bufferSize),
	)

	var readerWg sync.WaitGroup
	var workerWg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go p.startWorker(ctx, &workerWg, tradeChan, errChan, batchSize, tickerSeconds)
	}

	now := time.Now()

	for _, file := range files {
		readerWg.Add(1)
		go func(f string) {
			defer readerWg.Done()
			if err := ParseTXTFile(ctx, f, tradeChan, now); err != nil {
				log.L().Error("Error processing file", log.String("file", f), log.Error(err))
				select {
				case errChan <- err:
				default:
				}
			}
		}(file)
	}

	go func() {
		readerWg.Wait()
		close(tradeChan)
	}()

	workerWg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (p Processor) startWorker(
	ctx context.Context,
	wg *sync.WaitGroup,
	tradeChan <-chan db.CreateTradeParams,
	errChan chan<- error,
	batchSize int,
	tickerSeconds int,
) {
	defer wg.Done()

	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		select {
		case errChan <- err:
		default:
		}
		return
	}
	defer conn.Release()

	batch := make([]db.CreateTradeParams, 0, batchSize)
	ticker := time.NewTicker(time.Duration(tickerSeconds) * time.Second)
	defer ticker.Stop()

	flush := func() {
		if len(batch) > 0 {
			if err := p.repo.SaveBatchCopyFrom(ctx, conn.Conn(), batch); err != nil {
				log.L().Error("Error during batch insert", log.Error(err))
				select {
				case errChan <- err:
				default:
				}
			}
			batch = batch[:0]
		}
	}

	for {
		select {
		case trade, ok := <-tradeChan:
			if !ok {
				flush()
				return
			}
			batch = append(batch, trade)
			if len(batch) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			flush()
			return
		}
	}
}
