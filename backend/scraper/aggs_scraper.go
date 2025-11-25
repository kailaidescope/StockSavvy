package scraper

import (
	"encoding/json"
	"errors"
	"financial-helper/mongodb"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ScrapeTickerAggregatesOptions struct {
	collectionWindow *time.Duration
	collectionLimit  *int
}

type aggregatesOptionsJSON struct {
	CollectionWindow *int `json:"collection_window"`
	CollectionLimit  *int `json:"collection_limit"`
}
type aggregatesInstructionsJSON struct {
	Tickers   []string               `json:"tickers"`
	StartTime string                 `json:"start_time"`
	EndTime   string                 `json:"end_time"`
	Options   *aggregatesOptionsJSON `json:"options"`
}

// Reads scraping instructions from file and runs a scrape if instructions are valid
func (scraper *Scraper) ScrapeTickersAggregatesFromJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Join(errors.New("failed to read instructions file"), err)
	}

	var inst aggregatesInstructionsJSON
	if err := json.Unmarshal(data, &inst); err != nil {
		return errors.Join(errors.New("failed to parse instructions JSON"), err)
	}

	if len(inst.Tickers) == 0 {
		return errors.New("no tickers provided in JSON")
	}

	start, err := time.Parse("2006-01-02", inst.StartTime)
	if err != nil {
		return errors.Join(errors.New("invalid start_time"), err)
	}
	end, err := time.Parse("2006-01-02", inst.EndTime)
	if err != nil {
		return errors.Join(errors.New("invalid end_time"), err)
	}

	var opts ScrapeTickerAggregatesOptions
	if inst.Options != nil {
		if inst.Options.CollectionWindow != nil {
			d := time.Duration(*inst.Options.CollectionWindow) * 24 * time.Hour
			opts.collectionWindow = &d
		}
		if inst.Options.CollectionLimit != nil {
			l := *inst.Options.CollectionLimit
			opts.collectionLimit = &l
		}
	}

	return scraper.ScrapeTickersAggregates(inst.Tickers, start, end, opts)
}

func (scraper *Scraper) ScrapeTickersAggregates(symbols []string, start, end time.Time, options ScrapeTickerAggregatesOptions) error {
	if len(symbols) == 0 {
		return nil
	}

	total := len(symbols)

	// Stats
	startAll := time.Now()
	var tickersProcessed int
	var totalInsertedAggregates int
	totalSkippedAggregates := 0
	var totalTickerDuration time.Duration

	for i, symbol := range symbols {
		// compute elapsed & ETA based on completed tickers
		completed := tickersProcessed
		elapsed := time.Since(startAll)
		var eta time.Duration
		if completed > 0 {
			avgPer := elapsed / time.Duration(completed)
			eta = avgPer * time.Duration(total-completed)
		} else {
			eta = 0
		}

		// show which ticker is starting
		fmt.Printf("Starting (%d/%d): %s elapsed:%s eta:%s\n", i+1, total, symbol, formatDuration(elapsed), func() string {
			if eta > 0 {
				return formatDuration(eta)
			}
			return "??"
		}())

		tickerStart := time.Now()
		numInserted, numSkipped, err := scraper.ScrapeTickerAggregates(symbol, start, end, &options)
		duration := time.Since(tickerStart)

		// update stats
		tickersProcessed++
		totalInsertedAggregates += numInserted
		totalSkippedAggregates += numSkipped
		totalTickerDuration += duration

		// compute per-ticker progress summary
		percent := (float64(tickersProcessed) / float64(total)) * 100.0
		avgTimePerTicker := time.Duration(0)
		if tickersProcessed > 0 {
			avgTimePerTicker = totalTickerDuration / time.Duration(tickersProcessed)
		}
		etaPerTicker := time.Duration(0)
		if tickersProcessed > 0 {
			etaPerTicker = avgTimePerTicker * time.Duration(total-tickersProcessed)
		}
		avgAggregatesPerTicker := 0.0
		if tickersProcessed > 0 {
			avgAggregatesPerTicker = float64(totalInsertedAggregates) / float64(tickersProcessed)
		}
		avgSkippedPerTicker := 0.0
		if tickersProcessed > 0 {
			avgSkippedPerTicker = float64(totalSkippedAggregates) / float64(tickersProcessed)
		}

		// print progress summary after each ticker scrape
		fmt.Printf("\nPROGRESS: %.1f%%\neta=%s\navg_time_per_ticker=%s\naggregates_inserted=%d\nskipped_aggregates=%d\ntickers=%d\navg_aggregates_per_ticker=%.2f\navg_skipped_per_ticker=%.2f\n",
			percent, formatDuration(etaPerTicker), formatDuration(avgTimePerTicker), totalInsertedAggregates, totalSkippedAggregates, tickersProcessed, avgAggregatesPerTicker, avgSkippedPerTicker)

		if err != nil {
			// print partial results before returning
			totalElapsed := time.Since(startAll)
			fmt.Printf("\nRESULTS:\ntotal_time=%s\naverage_time_per_ticker=%s\ntickers_processed=%d\ninserted_aggregates=%d\nskipped_aggregates=%d\navg_aggregates_per_ticker=%.2f\navg_skipped_per_ticker=%.2f\n",
				formatDuration(totalElapsed), formatDuration(avgTimePerTicker), tickersProcessed, totalInsertedAggregates, totalSkippedAggregates, avgAggregatesPerTicker, avgSkippedPerTicker)
			return errors.Join(errors.New("error scraping "+symbol), err)
		}
	}

	// Final results log
	totalElapsed := time.Since(startAll)
	avgTimePerTicker := time.Duration(0)
	if tickersProcessed > 0 {
		avgTimePerTicker = totalTickerDuration / time.Duration(tickersProcessed)
	}
	avgAggregatesPerTicker := 0.0
	if tickersProcessed > 0 {
		avgAggregatesPerTicker = float64(totalInsertedAggregates) / float64(tickersProcessed)
	}
	avgSkippedPerTicker := 0.0
	if tickersProcessed > 0 {
		avgSkippedPerTicker = float64(totalSkippedAggregates) / float64(tickersProcessed)
	}
	fmt.Printf("\n\n=== RESULTS: ===\ntotal_time=%s\naverage_time_per_ticker=%s\ntickers_processed=%d\ninserted_aggregates=%d\nskipped_aggregates=%d\navg_aggregates_per_ticker=%.2f\navg_skipped_per_ticker=%.2f\n\n=============",
		formatDuration(totalElapsed), formatDuration(avgTimePerTicker), tickersProcessed, totalInsertedAggregates, totalSkippedAggregates, avgAggregatesPerTicker, avgSkippedPerTicker)
	return nil
}

func (scraper *Scraper) ScrapeTickerAggregates(symbol string, start, end time.Time, options *ScrapeTickerAggregatesOptions) (int, int, error) {
	// Validate input
	if start.After(end) {
		return 0, 0, errors.New("start time must be before end time")
	}

	// Set optional values & defaults
	collectionWindow := time.Hour * 24 * 7
	collectionLimit := 500
	if options != nil {
		if options.collectionWindow != nil {
			collectionWindow = *options.collectionWindow
		}
		if options.collectionLimit != nil {
			collectionLimit = *options.collectionLimit
		}
	}

	// Calculate number of steps (windows) for the progress bar based on days
	totalDays := int(math.Ceil(end.Sub(start).Hours() / 24.0))
	windowDays := int(math.Max(1, math.Floor(collectionWindow.Hours()/24.0)))
	numSteps := int(math.Ceil(float64(totalDays) / float64(windowDays)))
	if numSteps < 1 {
		numSteps = 1
	}

	bar := progressbar.NewOptions(numSteps,
		progressbar.OptionSetDescription(fmt.Sprintf("%s %s", symbol, start.Format(time.RFC3339))),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetPredictTime(false),
	)

	// Iterate through time series without overlapping days:
	// Each window covers exactly `windowDays` calendar days (inclusive).
	var insertedTotal int
	skippedTotal := 0
	stepDone := 0
	iterationStart := time.Now()

	for currentStart := start; !currentStart.After(end); {
		// Compute window end so that the window covers `windowDays` days inclusive.
		var currentEnd time.Time
		wd := windowDays
		if wd < 1 {
			wd = 1
		}
		// End is start + (wd-1) days (inclusive window of wd days)
		currentEnd = currentStart.Add(time.Duration(wd-1) * 24 * time.Hour)
		if currentEnd.After(end) {
			currentEnd = end
		}

		// compute elapsed & ETA for this ticker
		elapsed := time.Since(iterationStart)
		var eta time.Duration
		if stepDone > 0 {
			avgPerStep := elapsed / time.Duration(stepDone)
			eta = avgPerStep * time.Duration(numSteps-stepDone)
		} else {
			eta = 0
		}

		// run iteration in a closure so we can always advance the progress bar and the window
		func() {
			// update description with current start, elapsed & ETA
			bar.Describe(fmt.Sprintf("%s %s elapsed:%s eta:%s", symbol, currentStart.Format(time.RFC3339), formatDuration(elapsed), func() string {
				if eta > 0 {
					return formatDuration(eta)
				}
				return "??"
			}()))

			history, err := scraper.polygonClient.PolygonGetTickerHistory(symbol, currentStart, currentEnd, collectionLimit)
			if err != nil {
				// retry once
				history, err = scraper.polygonClient.PolygonGetTickerHistory(symbol, currentStart, currentEnd, collectionLimit)
				if err != nil {
					errLogger.Printf("Error receiving aggregates for %s from %s to %s : %s", symbol, currentStart.Format("2006-01-02T15:04:05Z"), currentEnd.Format("2006-01-02T15:04:05Z"), err.Error())
					return
				}
			}

			mongoAggs, err := mongodb.PolygonHistoryToAggs(*history)
			if err != nil {
				// retry once
				mongoAggs, err = mongodb.PolygonHistoryToAggs(*history)
				if err != nil {
					errLogger.Printf("Error converting to MongoDB aggregate types for %s from %s to %s : %s", symbol, currentStart.Format("2006-01-02T15:04:05Z"), currentEnd.Format("2006-01-02T15:04:05Z"), err.Error())
					return
				}
			}

			numInsertedAggs, err := mongodb.InsertAggregates(scraper.mongoClient, scraper.tickerDBName, mongoAggs)
			if err != nil {
				errLogger.Printf("Error inserting aggregates to MongoDB for %s from %s to %s : %s", symbol, currentStart.Format("2006-01-02T15:04:05Z"), currentEnd.Format("2006-01-02T15:04:05Z"), err.Error())
				return
			}
			if numInsertedAggs != len(mongoAggs) {
				errLogger.Printf("Some aggregates were not inserted for %s from %s to %s : %d/%d (inserted/total)", symbol, currentStart.Format("2006-01-02T15:04:05Z"), currentEnd.Format("2006-01-02T15:04:05Z"), numInsertedAggs, len(mongoAggs))
			}
			insertedTotal += numInsertedAggs
			skippedTotal += len(mongoAggs) - numInsertedAggs
		}()

		// advance progress and window regardless of success or failure
		_ = bar.Add(1)
		stepDone++

		// Move to the day after currentEnd to avoid overlap
		currentStart = currentEnd.Add(24 * time.Hour)
	}

	_ = bar.Finish()
	return insertedTotal, skippedTotal, nil
}
