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

	dbName := "test_stock_savvy"
	ticker := fmt.Sprintf("TEST-%d", time.Now().UnixNano())

	agg := TickerDailyAggregate{
		ID:           primitive.NewObjectID(),
		Ticker:       ticker,
		Volume:       12345.0,
		VWAP:         150.5,
		Open:         149.0,
		Close:        151.0,
		High:         152.0,
		Low:          148.5,
		Timestamp:    primitive.NewDateTimeFromTime(time.Now().UTC()),
		Transactions: 42,
		OTC:          false,
	}

	inserted, err := InsertHistory(testMongoClient, dbName, []TickerDailyAggregate{agg})
	if err != nil {
		t.Fatalf("InsertHistory returned error: %v", err)
	}
	if inserted != 1 {
		t.Fatalf("expected 1 inserted document, got %d", inserted)
	}

	// verify it exists in the DB
	coll := testMongoClient.Database(dbName).Collection("ticker_aggregates")
	var found TickerDailyAggregate
	if err := coll.FindOne(ctx, bson.M{"ticker": ticker, "timestamp": agg.Timestamp}).Decode(&found); err != nil {
		t.Fatalf("expected to find inserted aggregate, but got error: %v", err)
	}
	if found.Ticker != ticker {
		t.Fatalf("found aggregate ticker mismatch: expected %s got %s", ticker, found.Ticker)
	}

	/* // optional cleanup
	   if _, err := coll.DeleteMany(ctx, bson.M{"ticker": ticker}); err != nil {
	       t.Logf("cleanup DeleteMany error (non-fatal): %v", err)
	   } */
}

func TestInsertMultipleHistory(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	dbName := "test_stock_savvy"
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

	inserted, err := InsertHistory(testMongoClient, dbName, aggs)
	if err != nil {
		t.Fatalf("InsertHistory (batch) returned error: %v", err)
	}
	if inserted != 20 {
		t.Fatalf("expected 20 inserted documents, got %d", inserted)
	}

	// verify each inserted doc exists
	coll := testMongoClient.Database(dbName).Collection("ticker_aggregates")
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
	/* if _, err := coll.DeleteMany(ctx, bson.M{"ticker": bson.M{"$regex": "^" + prefix}}); err != nil {
	    t.Logf("cleanup DeleteMany error (non-fatal): %v", err)
	} */
}
