package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"negotiation-history-B3/internal/app"
	"negotiation-history-B3/internal/domain/trade"
)

func main() {
	cliApp := &cli.App{
		Name:  "b3-processor",
		Usage: "Processes B3 CSV files and allows data queries via CLI",
		Commands: []*cli.Command{
			{
				Name:  "process",
				Usage: "Processes one or more B3 CSV files",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "Path(s) to CSV files to process",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					files := c.StringSlice("file")

					return app.StartApp(context.Background(), func(tp trade.Processor) error {
						fmt.Println("Processing files:", files)
						if err := tp.ProcessFiles(files); err != nil {
							log.Printf("Error processing files: %v", err)
							return err
						}
						return nil
					})
				},
			},
			{
				Name:  "query",
				Usage: "Query aggregated data by ticker and optional start date",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ticker",
						Usage:    "Ticker symbol to query",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "date",
						Usage: "Optional start date (format: YYYY-MM-DD)",
					},
				},
				Action: func(c *cli.Context) error {
					ticker := c.String("ticker")
					dateStr := c.String("date")

					var fromDate *time.Time
					if dateStr != "" {
						loc, err := time.LoadLocation("America/Sao_Paulo")
						if err != nil {
							return fmt.Errorf("failed to load timezone: %w", err)
						}
						parsedDate, err := time.ParseInLocation("2006-01-02", dateStr, loc)
						if err != nil {
							return fmt.Errorf("invalid date format: %v", err)
						}
						fromDate = &parsedDate
					}

					return app.StartApp(context.Background(), func(fn trade.Finder) error {
						fmt.Printf("Querying data for ticker: %s, starting from: %v\n", ticker, fromDate)

						summary, err := fn.GetSummary(context.Background(), ticker, fromDate)
						if err != nil {
							log.Printf("Query error: %v", err)
							return err
						}

						fmt.Printf("Summary Trade Result:\n Ticker: %s\n Max Range Value: %.2f\n Max Daily Volume: %d\n",
							summary.Ticker, summary.MaxRangeValue, summary.MaxDailyVolume)
						return nil
					})
				},
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
