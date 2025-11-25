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

const DB_NAME string = "test_stock_savvy"

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

	inserted, err := InsertArticles(testMongoClient, DB_NAME, []Article{article})
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != 1 {
		t.Fatalf("expected 1 inserted document, got %d", inserted)
	}

	// verify it exists in the DB
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
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

	inserted, err := InsertArticles(testMongoClient, DB_NAME, articles)
	if err != nil {
		t.Fatalf("InsertArticles (batch) returned error: %v", err)
	}
	if inserted != 20 {
		t.Fatalf("expected 20 inserted documents, got %d", inserted)
	}

	// verify each inserted doc exists
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
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

// TestGetArticlesByTicker verifies GetArticlesByTicker returns articles containing the requested ticker.
func TestGetArticlesByTicker(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("test-query-ticker-%d-", time.Now().UnixNano())
	ticker := "TTICK"

	// create 3 articles, 2 containing ticker "TTICK" and 1 with a different ticker
	articles := make([]Article, 0, 3)
	for i := 0; i < 3; i++ {
		pid := fmt.Sprintf("%s%d-%d", prefix, i, rand.Intn(1_000_000))
		tks := []string{"OTHER"}
		if i < 2 {
			tks = []string{ticker}
		}
		a := Article{
			ID:          primitive.NewObjectID(),
			PolygonID:   pid,
			Publisher:   ArticlePublisher{Name: "Query Publisher"},
			Title:       fmt.Sprintf("Query Article %d", i),
			Author:      "query-tester",
			PublishedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
			ArticleURL:  "https://example.test/article",
			Tickers:     tks,
			ImageURL:    "https://example.test/image.png",
			Description: "Query test article",
			Keywords:    []string{"query", "test"},
		}
		articles = append(articles, a)
	}

	inserted, err := InsertArticles(testMongoClient, DB_NAME, articles)
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != len(articles) {
		t.Fatalf("expected %d inserted documents, got %d", len(articles), inserted)
	}

	// call function under test
	got, err := GetArticlesByTicker(testMongoClient, DB_NAME, ticker, 10, 0, 0)
	if err != nil {
		t.Fatalf("GetArticlesByTicker returned error: %v", err)
	}
	if len(got) < 2 {
		t.Fatalf("expected at least 2 articles for ticker %s, got %d", ticker, len(got))
	}

	//log.Printf("Got article(s): %#v", got)

	// cleanup
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

// TestGetArticlesByTickerOverRange ensures filtering by ticker + time range works.
func TestGetArticlesByTickerOverRange(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("test-query-range-%d-", time.Now().UnixNano())
	ticker := "RANGE"

	now := time.Now().UTC()
	inRangeTime := now.Add(-2 * time.Hour)
	outRangeTime := now.Add(-48 * time.Hour)

	a1 := Article{
		ID:          primitive.NewObjectID(),
		PolygonID:   fmt.Sprintf("%s1-%d", prefix, rand.Intn(1_000_000)),
		Publisher:   ArticlePublisher{Name: "Range Publisher"},
		Title:       "InRange Article",
		Author:      "range-tester",
		PublishedAt: primitive.NewDateTimeFromTime(inRangeTime),
		ArticleURL:  "https://example.test/article",
		Tickers:     []string{ticker},
	}
	a2 := Article{
		ID:          primitive.NewObjectID(),
		PolygonID:   fmt.Sprintf("%s2-%d", prefix, rand.Intn(1_000_000)),
		Publisher:   ArticlePublisher{Name: "Range Publisher"},
		Title:       "OutRange Article",
		Author:      "range-tester",
		PublishedAt: primitive.NewDateTimeFromTime(outRangeTime),
		ArticleURL:  "https://example.test/article",
		Tickers:     []string{ticker},
	}

	inserted, err := InsertArticles(testMongoClient, DB_NAME, []Article{a1, a2})
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != 2 {
		t.Fatalf("expected 2 inserted documents, got %d", inserted)
	}

	start := now.Add(-24 * time.Hour)
	end := now

	got, err := GetArticlesByTickerOverRange(testMongoClient, DB_NAME, ticker, start, end, 10, 0, 0)
	if err != nil {
		t.Fatalf("GetArticlesByTickerOverRange returned error: %v", err)
	}
	// should include only the in-range article
	if len(got) != 1 {
		t.Fatalf("expected 1 article in range, got %d", len(got))
	}
	if got[0].PolygonID != a1.PolygonID {
		t.Fatalf("expected polygon_id %s, got %s", a1.PolygonID, got[0].PolygonID)
	}

	//log.printf("Got article(s): %#v", got)

	// cleanup
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

// TestGetArticlesOverRange verifies GetArticlesOverRange returns articles filtered by published_at range.
func TestGetArticlesOverRange(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("test-query-over-%d-", time.Now().UnixNano())

	now := time.Now().UTC()
	inside := now.Add(-3 * time.Hour)
	inside2 := now.Add(-2 * time.Hour)
	outside := now.Add(-72 * time.Hour)

	a1 := Article{
		ID:          primitive.NewObjectID(),
		PolygonID:   fmt.Sprintf("%s1-%d", prefix, rand.Intn(1_000_000)),
		Publisher:   ArticlePublisher{Name: "Over Publisher"},
		Title:       "Inside 1",
		Author:      "over-tester",
		PublishedAt: primitive.NewDateTimeFromTime(inside),
		ArticleURL:  "https://example.test/article",
		Tickers:     []string{"X"},
	}
	a2 := Article{
		ID:          primitive.NewObjectID(),
		PolygonID:   fmt.Sprintf("%s2-%d", prefix, rand.Intn(1_000_000)),
		Publisher:   ArticlePublisher{Name: "Over Publisher"},
		Title:       "Inside 2",
		Author:      "over-tester",
		PublishedAt: primitive.NewDateTimeFromTime(inside2),
		ArticleURL:  "https://example.test/article",
		Tickers:     []string{"Y"},
	}
	a3 := Article{
		ID:          primitive.NewObjectID(),
		PolygonID:   fmt.Sprintf("%s3-%d", prefix, rand.Intn(1_000_000)),
		Publisher:   ArticlePublisher{Name: "Over Publisher"},
		Title:       "Outside",
		Author:      "over-tester",
		PublishedAt: primitive.NewDateTimeFromTime(outside),
		ArticleURL:  "https://example.test/article",
		Tickers:     []string{"Z"},
	}

	inserted, err := InsertArticles(testMongoClient, DB_NAME, []Article{a1, a2, a3})
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != 3 {
		t.Fatalf("expected 3 inserted documents, got %d", inserted)
	}

	start := now.Add(-24 * time.Hour)
	end := now

	got, err := GetArticlesOverRange(testMongoClient, DB_NAME, start, end, 10, 0, 0)
	if err != nil {
		t.Fatalf("GetArticlesOverRange returned error: %v", err)
	}

	// we expect at least the two inside articles to be returned (there may be other data in DB)
	foundIDs := map[string]struct{}{}
	for _, g := range got {
		foundIDs[g.PolygonID] = struct{}{}
	}
	if _, ok := foundIDs[a1.PolygonID]; !ok {
		t.Fatalf("expected to find article %s in results", a1.PolygonID)
	}
	if _, ok := foundIDs[a2.PolygonID]; !ok {
		t.Fatalf("expected to find article %s in results", a2.PolygonID)
	}

	//log.printf("Got article(s): %#v", got)

	// cleanup
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

// TestGetArticlesByTickerPagination verifies page-based pagination for GetArticlesByTicker.
func TestGetArticlesByTickerPagination(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("test-pag-ticker-%d-", time.Now().UnixNano())
	ticker := "PAGTICK"
	count := 7
	pageSize := 3

	// create `count` articles with descending published_at (i=0 newest)
	articles := make([]Article, 0, count)
	expectedOrder := make([]string, 0, count)
	now := time.Now().UTC()
	for i := 0; i < count; i++ {
		pid := fmt.Sprintf("%s%d-%d", prefix, i, rand.Intn(1_000_000))
		pub := now.Add(-time.Duration(i) * time.Minute) // i=0 newest
		a := Article{
			ID:          primitive.NewObjectID(),
			PolygonID:   pid,
			Publisher:   ArticlePublisher{Name: "Pag Publisher"},
			Title:       fmt.Sprintf("Pag Article %d", i),
			Author:      "pag-tester",
			PublishedAt: primitive.NewDateTimeFromTime(pub),
			ArticleURL:  "https://example.test/article",
			Tickers:     []string{ticker},
		}
		articles = append(articles, a)
		expectedOrder = append(expectedOrder, pid)
	}

	inserted, err := InsertArticles(testMongoClient, DB_NAME, articles)
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != len(articles) {
		t.Fatalf("expected %d inserted documents, got %d", len(articles), inserted)
	}

	// page 1 should return newest 3 (expectedOrder[0..2])
	got1, err := GetArticlesByTicker(testMongoClient, DB_NAME, ticker, 0, 1, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesByTicker page1 returned error: %v", err)
	}
	if len(got1) != pageSize {
		t.Fatalf("expected %d results on page1, got %d", pageSize, len(got1))
	}
	for i := 0; i < pageSize; i++ {
		if got1[i].PolygonID != expectedOrder[i] {
			t.Fatalf("page1: expected polygon_id %s at index %d, got %s", expectedOrder[i], i, got1[i].PolygonID)
		}
	}

	// page 2 should return next 3 (expectedOrder[3..5])
	got2, err := GetArticlesByTicker(testMongoClient, DB_NAME, ticker, 0, 2, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesByTicker page2 returned error: %v", err)
	}
	if len(got2) != pageSize {
		t.Fatalf("expected %d results on page2, got %d", pageSize, len(got2))
	}
	for i := 0; i < pageSize; i++ {
		if got2[i].PolygonID != expectedOrder[i+pageSize] {
			t.Fatalf("page2: expected polygon_id %s at index %d, got %s", expectedOrder[i+pageSize], i, got2[i].PolygonID)
		}
	}

	// page 3 should return the remainder (1 item)
	got3, err := GetArticlesByTicker(testMongoClient, DB_NAME, ticker, 0, 3, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesByTicker page3 returned error: %v", err)
	}
	expectedLast := count - 2*pageSize // 7 - 6 = 1
	if len(got3) != expectedLast {
		t.Fatalf("expected %d results on page3, got %d", expectedLast, len(got3))
	}
	if got3[0].PolygonID != expectedOrder[2*pageSize] {
		t.Fatalf("page3: expected polygon_id %s, got %s", expectedOrder[2*pageSize], got3[0].PolygonID)
	}

	// cleanup
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

// TestGetArticlesByTickerOverRangePagination verifies pagination for ticker+range queries.
func TestGetArticlesByTickerOverRangePagination(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("test-pag-range-%d-", time.Now().UnixNano())
	ticker := "RANGEPAG"
	total := 8
	pageSize := 3

	now := time.Now().UTC()
	articles := make([]Article, 0, total)
	expectedOrder := make([]string, 0, total)

	// create `total` articles across 8 minutes so ordering is clear
	for i := 0; i < total; i++ {
		pid := fmt.Sprintf("%s%d-%d", prefix, i, rand.Intn(1_000_000))
		pub := now.Add(-time.Duration(i) * time.Minute) // i=0 newest
		a := Article{
			ID:          primitive.NewObjectID(),
			PolygonID:   pid,
			Publisher:   ArticlePublisher{Name: "RangePag Publisher"},
			Title:       fmt.Sprintf("RangePag Article %d", i),
			Author:      "rangepag-tester",
			PublishedAt: primitive.NewDateTimeFromTime(pub),
			ArticleURL:  "https://example.test/article",
			Tickers:     []string{ticker},
		}
		articles = append(articles, a)
		expectedOrder = append(expectedOrder, pid)
	}

	inserted, err := InsertArticles(testMongoClient, DB_NAME, articles)
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != len(articles) {
		t.Fatalf("expected %d inserted documents, got %d", len(articles), inserted)
	}

	// choose a start/end that include all created articles (last 24 hours)
	start := now.Add(-24 * time.Hour)
	end := now.Add(1 * time.Minute)

	// page 1
	got1, err := GetArticlesByTickerOverRange(testMongoClient, DB_NAME, ticker, start, end, 0, 1, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesByTickerOverRange page1 error: %v", err)
	}
	if len(got1) != pageSize {
		t.Fatalf("expected %d results on page1, got %d", pageSize, len(got1))
	}
	for i := 0; i < pageSize; i++ {
		if got1[i].PolygonID != expectedOrder[i] {
			t.Fatalf("page1: expected polygon_id %s at index %d, got %s", expectedOrder[i], i, got1[i].PolygonID)
		}
	}

	// page 2
	got2, err := GetArticlesByTickerOverRange(testMongoClient, DB_NAME, ticker, start, end, 0, 2, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesByTickerOverRange page2 error: %v", err)
	}
	if len(got2) != pageSize {
		t.Fatalf("expected %d results on page2, got %d", pageSize, len(got2))
	}
	for i := 0; i < pageSize; i++ {
		if got2[i].PolygonID != expectedOrder[i+pageSize] {
			t.Fatalf("page2: expected polygon_id %s at index %d, got %s", expectedOrder[i+pageSize], i, got2[i].PolygonID)
		}
	}

	// page 3 (should have total - 6 = 2 items)
	got3, err := GetArticlesByTickerOverRange(testMongoClient, DB_NAME, ticker, start, end, 0, 3, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesByTickerOverRange page3 error: %v", err)
	}
	expectedLast := total - 2*pageSize
	if len(got3) != expectedLast {
		t.Fatalf("expected %d results on page3, got %d", expectedLast, len(got3))
	}
	if got3[0].PolygonID != expectedOrder[2*pageSize] {
		t.Fatalf("page3: expected polygon_id %s, got %s", expectedOrder[2*pageSize], got3[0].PolygonID)
	}

	// cleanup
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}

// TestGetArticlesOverRangePagination verifies pagination for global range queries.
func TestGetArticlesOverRangePagination(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("test mongo client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("test-pag-over-%d-", time.Now().UnixNano())
	total := 9
	pageSize := 4

	now := time.Now().UTC()
	articles := make([]Article, 0, total)
	expectedOrder := make([]string, 0, total)

	// create `total` articles (different tickers) so global range covers them
	for i := 0; i < total; i++ {
		pid := fmt.Sprintf("%s%d-%d", prefix, i, rand.Intn(1_000_000))
		pub := now.Add(-time.Duration(i) * time.Minute)
		a := Article{
			ID:          primitive.NewObjectID(),
			PolygonID:   pid,
			Publisher:   ArticlePublisher{Name: "OverPag Publisher"},
			Title:       fmt.Sprintf("OverPag Article %d", i),
			Author:      "overpag-tester",
			PublishedAt: primitive.NewDateTimeFromTime(pub),
			ArticleURL:  "https://example.test/article",
			Tickers:     []string{fmt.Sprintf("T%d", i)},
		}
		articles = append(articles, a)
		expectedOrder = append(expectedOrder, pid)
	}

	inserted, err := InsertArticles(testMongoClient, DB_NAME, articles)
	if err != nil {
		t.Fatalf("InsertArticles returned error: %v", err)
	}
	if inserted != len(articles) {
		t.Fatalf("expected %d inserted documents, got %d", len(articles), inserted)
	}

	start := now.Add(-24 * time.Hour)
	end := now.Add(1 * time.Minute)

	// page 1
	got1, err := GetArticlesOverRange(testMongoClient, DB_NAME, start, end, 0, 1, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesOverRange page1 error: %v", err)
	}
	if len(got1) != pageSize {
		t.Fatalf("expected %d results on page1, got %d", pageSize, len(got1))
	}
	for i := 0; i < pageSize; i++ {
		if got1[i].PolygonID != expectedOrder[i] {
			t.Fatalf("page1: expected polygon_id %s at index %d, got %s", expectedOrder[i], i, got1[i].PolygonID)
		}
	}

	// page 2
	got2, err := GetArticlesOverRange(testMongoClient, DB_NAME, start, end, 0, 2, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesOverRange page2 error: %v", err)
	}
	if len(got2) != pageSize {
		t.Fatalf("expected %d results on page2, got %d", pageSize, len(got2))
	}
	for i := 0; i < pageSize; i++ {
		if got2[i].PolygonID != expectedOrder[i+pageSize] {
			t.Fatalf("page2: expected polygon_id %s at index %d, got %s", expectedOrder[i+pageSize], i, got2[i].PolygonID)
		}
	}

	// page 3 (should have total - 8 = 1)
	got3, err := GetArticlesOverRange(testMongoClient, DB_NAME, start, end, 0, 3, pageSize)
	if err != nil {
		t.Fatalf("GetArticlesOverRange page3 error: %v", err)
	}
	expectedLast := total - 2*pageSize
	if len(got3) != expectedLast {
		t.Fatalf("expected %d results on page3, got %d", expectedLast, len(got3))
	}
	if got3[0].PolygonID != expectedOrder[2*pageSize] {
		t.Fatalf("page3: expected polygon_id %s, got %s", expectedOrder[2*pageSize], got3[0].PolygonID)
	}

	// cleanup
	coll := testMongoClient.Database(DB_NAME).Collection("ticker_news")
	if _, err := coll.DeleteMany(ctx, bson.M{"polygon_id": bson.M{"$regex": "^" + prefix}}); err != nil {
		t.Logf("cleanup error: %v", err)
	}
}
