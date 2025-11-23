package mongodb

import (
	"context"
	"errors"
	"financial-helper/environment"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var testMongoClient *mongo.Client
var testTeardown func()
var testInitErr error

func TestMain(m *testing.M) {
	// Initialize mongo client once for all tests
	testMongoClient, testTeardown, testInitErr = func() (*mongo.Client, func(), error) {
		if err := environment.LoadEnvironment(); err != nil {
			return nil, nil, errors.Join(errors.New("failed to load environment for testing"), err)
		}
		vars, _, err := environment.LoadVars()
		if err != nil {
			return nil, nil, errors.Join(errors.New("failed to load env variables"), err)
		}

		mongoPort, _ := strconv.Atoi(vars["MONGO_PORT"])
		client, err := GetMongoDBInstance(vars["MONGO_INITDB_ROOT_USERNAME"], vars["MONGO_INITDB_ROOT_PASSWORD"], vars["MONGO_HOST"], mongoPort)
		if err != nil {
			return nil, nil, err
		}

		return client, func() {
			log.Println("Tearing down test environment...")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = client.Disconnect(ctx)
		}, nil
	}()

	if testInitErr != nil {
		log.Printf("failed to initialize test mongo client: %v\n", testInitErr)
		os.Exit(1)
	}

	code := m.Run()

	if testTeardown != nil {
		testTeardown()
	}
	os.Exit(code)
}

// TestInsertArticles creates a stub Article and attempts to insert it into the test DB.
func TestInsertArticles(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dbName := "test_stock_savvy"
	polygonID := fmt.Sprintf("test-article-%d", time.Now().UnixNano())

	article := Article{
		ID:          primitive.NewObjectID(),
		PolygonID:   polygonID,
		Publisher:   ArticlePublisher{Name: "Test Publisher", HomepageURL: "https://example.test"},
		Title:       "Test Article Title",
		Author:      "unit-tester",
		PublishedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
		ArticleURL:  "https://example.test/article",
		Tickers:     []string{"AAPL"},
		ImageURL:    "https://example.test/image.png",
		Description: "This is a test article inserted by unit tests.",
		Keywords:    []string{"test", "unit"},
		Insights:    []ArticleInsight{{Ticker: "AAPL", Sentiment: "neutral", SentimentReasoning: "stub"}},
	}

	inserted, err := InsertArticles(testMongoClient, dbName, []Article{article})
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != 1 {
		t.Fatalf("expected 1 inserted document, got %d", inserted)
	}

	// verify it exists in the DB
	coll := testMongoClient.Database(dbName).Collection("ticker_news")
	var found Article
	if err := coll.FindOne(ctx, bson.M{"polygon_id": polygonID}).Decode(&found); err != nil {
		t.Fatalf("expected to find inserted article, but got error: %v", err)
	}
	if found.PolygonID != polygonID {
		t.Fatalf("found article polygon_id mismatch: expected %s got %s", polygonID, found.PolygonID)
	}

	// cleanup
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": polygonID}); err != nil {
		t.Logf("cleanup DeleteMany error (non-fatal): %v", err)
	}
}

// TestInsertMultipleArticles inserts a batch of 20 articles with randomized polygon_id and published_at.
func TestInsertMultipleArticles(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	dbName := "test_stock_savvy"
	prefix := fmt.Sprintf("test-article-batch-%d-", time.Now().UnixNano())

	tickers := []string{"AAPL", "MSFT", "GOOG", "TSLA", "AMZN"}

	articles := make([]Article, 0, 20)
	polygonIDs := make([]string, 0, 20)

	for i := 0; i < 20; i++ {
		pid := fmt.Sprintf("%s%d-%d", prefix, i, rand.Intn(1_000_000))
		polygonIDs = append(polygonIDs, pid)

		// published_at random within last 30 days
		secsBack := rand.Intn(30 * 24 * 3600)
		pubTime := time.Now().UTC().Add(-time.Duration(secsBack) * time.Second)

		a := Article{
			ID:          primitive.NewObjectID(),
			PolygonID:   pid,
			Publisher:   ArticlePublisher{Name: "Batch Publisher", HomepageURL: "https://example.test"},
			Title:       fmt.Sprintf("Batch Article %d", i),
			Author:      "batch-tester",
			PublishedAt: primitive.NewDateTimeFromTime(pubTime),
			ArticleURL:  "https://example.test/article",
			Tickers:     []string{tickers[rand.Intn(len(tickers))]},
			ImageURL:    "https://example.test/image.png",
			Description: "Batch insert test article",
			Keywords:    []string{"batch", "test"},
			Insights:    []ArticleInsight{{Ticker: "AAPL", Sentiment: "neutral", SentimentReasoning: "batch"}},
		}
		articles = append(articles, a)
	}

	inserted, err := InsertArticles(testMongoClient, dbName, articles)
	if err != nil {
		t.Fatalf("InsertArticles (batch) returned error: %v", err)
	}
	if inserted != 20 {
		t.Fatalf("expected 20 inserted documents, got %d", inserted)
	}

	// verify each inserted doc exists
	coll := testMongoClient.Database(dbName).Collection("ticker_news")
	for _, pid := range polygonIDs {
		var found Article
		if err := coll.FindOne(ctx, bson.M{"polygon_id": pid}).Decode(&found); err != nil {
			t.Fatalf("expected to find inserted article %s, but got error: %v", pid, err)
		}
		if found.PolygonID != pid {
			t.Fatalf("found article polygon_id mismatch: expected %s got %s", pid, found.PolygonID)
		}
	}

	// cleanup - remove the batch by prefix
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup DeleteMany error (non-fatal): %v", err)
	}
}
