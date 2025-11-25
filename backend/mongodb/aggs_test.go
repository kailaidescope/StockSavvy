package mongodb

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertHistory(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ticker := fmt.Sprintf("TEST-%d", time.Now().UnixNano())

	agg := TickerDailyAggregate{
		ID:           primitive.NewObjectID(),
		Ticker:       ticker,
		Volume:       8189.0,
		VWAP:         150.5,
		Open:         149.0,
		Close:        151.0,
		High:         152.0,
		Low:          148.5,
		Timestamp:    primitive.NewDateTimeFromTime(time.Now().UTC()),
		Transactions: 42,
		OTC:          false,
	}

	/* tsStr := "2025-11-03T16:46:51.477+00:00"
	tParsed, err := time.Parse(time.RFC3339Nano, tsStr)
	if err != nil {
		t.Fatalf("failed to parse timestamp %q: %v", tsStr, err)
	}
	timestamp := primitive.NewDateTimeFromTime(tParsed.UTC())

	ticker = "batch-test-1764109200477013000-9-TSLA"

	id, err := primitive.ObjectIDFromHex("69262b90fed140f33f67253b")
	if err != nil {
		t.Fatalf("Faield to parse obj id %s", err.Error())
	}

	agg.Timestamp = timestamp
	agg.ID = id
	agg.Ticker = ticker */

	inserted, err := InsertAggregates(testMongoClient, DB_NAME, []TickerDailyAggregate{agg})
	if err != nil {
		t.Fatalf("InsertHistory returned error: %v", err)
	}
	if inserted != 1 {
		t.Fatalf("expected 1 inserted document, got %d", inserted)
	}

	// verify it exists in the DB
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_aggregates")
	var found TickerDailyAggregate
	if err := coll.FindOne(ctx, bson.M{"ticker": ticker, "timestamp": agg.Timestamp}).Decode(&found); err != nil {
		t.Fatalf("expected to find inserted aggregate, but got error: %v", err)
	}
	if found.Ticker != ticker {
		t.Fatalf("found aggregate ticker mismatch: expected %s got %s", ticker, found.Ticker)
	}

	// optional cleanup
	if _, err := coll.DeleteMany(ctx, bson.M{"ticker": ticker}); err != nil {
		t.Logf("cleanup DeleteMany error (non-fatal): %v", err)
	}
}

func TestInsertMultipleHistory(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("batch-test-%d-", time.Now().UnixNano())

	aggs := make([]TickerDailyAggregate, 0, 20)
	tickers := []string{"AAPL", "MSFT", "GOOG", "TSLA", "AMZN"}
	tickerNames := make([]string, 0, 20)

	for i := 0; i < 20; i++ {
		ticker := fmt.Sprintf("%s%d-%s", prefix, i, tickers[rand.Intn(5)])
		tickerNames = append(tickerNames, ticker)

		// timestamp random within last 30 days
		secsBack := rand.Intn(30 * 24 * 3600)
		ts := time.Now().UTC().Add(-time.Duration(secsBack) * time.Second)

		a := TickerDailyAggregate{
			ID:           primitive.NewObjectID(),
			Ticker:       ticker,
			Volume:       float64(1000 + rand.Intn(100000)),
			VWAP:         100 + rand.Float64()*50,
			Open:         100 + rand.Float64()*50,
			Close:        100 + rand.Float64()*50,
			High:         150 + rand.Float64()*10,
			Low:          90 + rand.Float64()*10,
			Timestamp:    primitive.NewDateTimeFromTime(ts),
			Transactions: rand.Intn(1000),
			OTC:          false,
		}
		aggs = append(aggs, a)
	}

	inserted, err := InsertAggregates(testMongoClient, DB_NAME, aggs)
	if err != nil {
		t.Fatalf("InsertHistory (batch) returned error: %v", err)
	}
	if inserted != 20 {
		t.Fatalf("expected 20 inserted documents, got %d", inserted)
	}

	// verify each inserted doc exists
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_aggregates")
	for _, tk := range tickerNames {
		var found TickerDailyAggregate
		if err := coll.FindOne(ctx, bson.M{"ticker": tk}).Decode(&found); err != nil {
			t.Fatalf("expected to find inserted aggregate %s, but got error: %v", tk, err)
		}
		if found.Ticker != tk {
			t.Fatalf("found aggregate ticker mismatch: expected %s got %s", tk, found.Ticker)
		}
	}

	// cleanup - optional
	if _, err := coll.DeleteMany(ctx, bson.M{"ticker": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup DeleteMany error (non-fatal): %v", err)
	}
}

func makeAgg(ticker string, ts time.Time, seq int) TickerDailyAggregate {
	return TickerDailyAggregate{
		ID:           primitive.NewObjectID(),
		Ticker:       ticker,
		Volume:       float64(1000 + seq),
		VWAP:         100 + float64(seq),
		Open:         100 + float64(seq),
		Close:        100 + float64(seq),
		High:         150 + float64(seq),
		Low:          90 + float64(seq),
		Timestamp:    primitive.NewDateTimeFromTime(ts),
		Transactions: seq,
		OTC:          false,
	}
}

func TestGetAggregatesByTicker_NoPagination(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ticker := fmt.Sprintf("QP-%d-NP", time.Now().UnixNano())
	// create 5 daily aggregates with increasing timestamps (older -> newer)
	now := time.Now().UTC()
	aggs := make([]TickerDailyAggregate, 0, 5)
	for i := 0; i < 5; i++ {
		ts := now.Add(time.Duration(i-4) * 24 * time.Hour) // oldest first
		aggs = append(aggs, makeAgg(ticker, ts, i))
	}

	// insert
	ins, err := InsertAggregates(testMongoClient, DB_NAME, aggs)
	if err != nil {
		t.Fatalf("InsertAggregates error: %v", err)
	}
	if ins != len(aggs) {
		t.Fatalf("expected %d inserted, got %d", len(aggs), ins)
	}

	// query without pagination (limit,page,pageSize = 0)
	got, err := GetAggregatesByTicker(testMongoClient, DB_NAME, ticker, 0, 0, 0)
	if err != nil {
		t.Fatalf("GetAggregatesByTicker error: %v", err)
	}
	if len(got) != len(aggs) {
		t.Fatalf("expected %d results, got %d", len(aggs), len(got))
	}

	// verify ascending order by timestamp
	for i := 1; i < len(got); i++ {
		if int64(got[i].Timestamp) <= int64(got[i-1].Timestamp) {
			t.Fatalf("results not strictly ascending by timestamp at pos %d", i)
		}
	}

	// cleanup
	if _, err := testMongoClient.Database(DB_NAME).Collection("ticker_aggregates").DeleteMany(ctx, bson.M{"ticker": ticker}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

func TestGetAggregatesByTicker_Pagination(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ticker := fmt.Sprintf("QP-%d-PG", time.Now().UnixNano())
	now := time.Now().UTC()
	total := 7
	aggs := make([]TickerDailyAggregate, 0, total)
	for i := 0; i < total; i++ {
		ts := now.Add(time.Duration(i-6) * 24 * time.Hour)
		aggs = append(aggs, makeAgg(ticker, ts, i))
	}

	ins, err := InsertAggregates(testMongoClient, DB_NAME, aggs)
	if err != nil {
		t.Fatalf("InsertAggregates error: %v", err)
	}
	if ins != len(aggs) {
		t.Fatalf("expected %d inserted, got %d", len(aggs), ins)
	}

	pageSize := 3
	// page 1
	got1, err := GetAggregatesByTicker(testMongoClient, DB_NAME, ticker, 0, 1, pageSize)
	if err != nil {
		t.Fatalf("GetAggregatesByTicker page1 error: %v", err)
	}
	if len(got1) != pageSize {
		t.Fatalf("expected page1 size %d, got %d", pageSize, len(got1))
	}
	// page 2
	got2, err := GetAggregatesByTicker(testMongoClient, DB_NAME, ticker, 0, 2, pageSize)
	if err != nil {
		t.Fatalf("GetAggregatesByTicker page2 error: %v", err)
	}
	if len(got2) != pageSize {
		t.Fatalf("expected page2 size %d, got %d", pageSize, len(got2))
	}
	// page 3 (should contain the remainder)
	got3, err := GetAggregatesByTicker(testMongoClient, DB_NAME, ticker, 0, 3, pageSize)
	if err != nil {
		t.Fatalf("GetAggregatesByTicker page3 error: %v", err)
	}
	expectedLast := total - 2*pageSize
	if len(got3) != expectedLast {
		t.Fatalf("expected page3 size %d, got %d", expectedLast, len(got3))
	}

	// verify continuity: last ts of page1 < first ts of page2, etc.
	if int64(got1[len(got1)-1].Timestamp) >= int64(got2[0].Timestamp) {
		t.Fatalf("page continuity broken between page1 and page2")
	}
	if int64(got2[len(got2)-1].Timestamp) >= int64(got3[0].Timestamp) {
		t.Fatalf("page continuity broken between page2 and page3")
	}

	// cleanup
	if _, err := testMongoClient.Database(DB_NAME).Collection("ticker_aggregates").DeleteMany(ctx, bson.M{"ticker": ticker}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

func TestGetAggregatesByTickerOverRange_NoPaginationAndPagination(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ticker := fmt.Sprintf("QPR-%d", time.Now().UnixNano())
	now := time.Now().UTC()
	// create 6 daily aggregates spanning 6 days
	aggs := make([]TickerDailyAggregate, 0, 6)
	for i := 0; i < 6; i++ {
		ts := now.Add(time.Duration(i-5) * 24 * time.Hour) // oldest first
		aggs = append(aggs, makeAgg(ticker, ts, i))
	}

	ins, err := InsertAggregates(testMongoClient, DB_NAME, aggs)
	if err != nil {
		t.Fatalf("InsertAggregates error: %v", err)
	}
	if ins != len(aggs) {
		t.Fatalf("expected %d inserted, got %d", len(aggs), ins)
	}

	// range: middle 3 days (indexes 1..3)
	start := aggs[1].Timestamp.Time()
	end := aggs[3].Timestamp.Time()

	// no pagination
	got, err := GetAggregatesByTickerOverRange(testMongoClient, DB_NAME, ticker, start, end, 0, 0, 0)
	if err != nil {
		t.Fatalf("GetAggregatesByTickerOverRange error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 results in range, got %d", len(got))
	}

	// with pagination pageSize=2 page=1 should return first 2 in that range
	gotP1, err := GetAggregatesByTickerOverRange(testMongoClient, DB_NAME, ticker, start, end, 0, 1, 2)
	if err != nil {
		t.Fatalf("GetAggregatesByTickerOverRange paged error: %v", err)
	}
	if len(gotP1) != 2 {
		t.Fatalf("expected paged size 2, got %d", len(gotP1))
	}
	// page 2 should return the remaining 1
	gotP2, err := GetAggregatesByTickerOverRange(testMongoClient, DB_NAME, ticker, start, end, 0, 2, 2)
	if err != nil {
		t.Fatalf("GetAggregatesByTickerOverRange paged page2 error: %v", err)
	}
	if len(gotP2) != 1 {
		t.Fatalf("expected paged size 1, got %d", len(gotP2))
	}

	// cleanup
	if _, err := testMongoClient.Database(DB_NAME).Collection("ticker_aggregates").DeleteMany(ctx, bson.M{"ticker": ticker}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

func TestGetAggregatesOverRange_NoPaginationAndPagination(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("QOR-%d-", time.Now().UnixNano())
	now := time.Now().UTC()

	// create aggregates across 3 tickers, 3 days each (9 docs total)
	allAggs := make([]TickerDailyAggregate, 0, 9)
	tickers := []string{prefix + "A", prefix + "B", prefix + "C"}
	for ti := 0; ti < len(tickers); ti++ {
		for di := 0; di < 3; di++ {
			ts := now.Add(time.Duration(di-2) * 24 * time.Hour) // oldest first per ticker but timestamps overlap across tickers
			allAggs = append(allAggs, makeAgg(tickers[ti], ts, ti*10+di))
		}
	}

	ins, err := InsertAggregates(testMongoClient, DB_NAME, allAggs)
	if err != nil {
		t.Fatalf("InsertAggregates error: %v", err)
	}
	if ins != len(allAggs) {
		t.Fatalf("expected %d inserted, got %d", len(allAggs), ins)
	}

	// query range to include only the middle day (di==1 across tickers)
	start := allAggs[1].Timestamp.Time()
	end := start

	// no pagination => should return 3 docs (one per ticker)
	got, err := GetAggregatesOverRange(testMongoClient, DB_NAME, start, end, 0, 0, 0)
	if err != nil {
		t.Fatalf("GetAggregatesOverRange error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 results in range, got %d", len(got))
	}

	// pagination pageSize=2 page=1 should return 2 docs, page=2 -> 1 doc
	gotP1, err := GetAggregatesOverRange(testMongoClient, DB_NAME, start, end, 0, 1, 2)
	if err != nil {
		t.Fatalf("GetAggregatesOverRange paged error: %v", err)
	}
	if len(gotP1) != 2 {
		t.Fatalf("expected paged size 2, got %d", len(gotP1))
	}
	gotP2, err := GetAggregatesOverRange(testMongoClient, DB_NAME, start, end, 0, 2, 2)
	if err != nil {
		t.Fatalf("GetAggregatesOverRange paged page2 error: %v", err)
	}
	if len(gotP2) != 1 {
		t.Fatalf("expected paged size 1, got %d", len(gotP2))
	}

	// cleanup
	if _, err := testMongoClient.Database(DB_NAME).Collection("ticker_aggregates").DeleteMany(ctx, bson.M{"ticker": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}
