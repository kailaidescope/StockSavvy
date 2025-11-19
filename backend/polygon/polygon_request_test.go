package polygon

import (
	"errors"
	"financial-helper/environment"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var testTicker string = "AAPL"

var polygonConnection *PolygonConnection
var testTeardown func()
var testInitErr error

func TestMain(m *testing.M) {
	// Initialize server once for all tests
	polygonConnection, testTeardown, testInitErr = func() (*PolygonConnection, func(), error) {
		if err := environment.LoadEnvironment(); err != nil {
			return nil, nil, errors.Join(errors.New("failed to load environment for testing"), err)
		}
		vars, polygonKeys, err := environment.LoadVars()
		if err != nil {
			return nil, nil, errors.Join(errors.New("failed to load env variables"), err)
		}
		throttleTimeInt, _ := strconv.Atoi(vars["THROTTLE_TIME"])
		server := GetPolygonConnection(polygonKeys, time.Duration(throttleTimeInt)*time.Second)
		if err != nil {
			return nil, nil, err
		}
		// Note: don't start the router unless you need it for integration tests
		return server, func() {
			log.Println("Tearing down test environment...")
		}, nil
	}()

	if testInitErr != nil {
		// Fail early; individual tests can still choose to skip if API key is missing.
		log.Printf("failed to initialize test server: %v\n", testInitErr)
		os.Exit(1)
	}

	code := m.Run()

	if testTeardown != nil {
		testTeardown()
	}
	os.Exit(code)
}

func TestPolygonGetTicker(t *testing.T) {
	if polygonConnection == nil {
		t.Skip("test server not initialized")
	}

	resp, err := polygonConnection.PolygonGetTicker(testTicker)
	if err != nil {
		t.Fatalf("PolygonGetTicker error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if resp.Count == nil || *resp.Count == 0 {
		t.Fatalf("expected count > 0")
	}
	if resp.Results == nil || len(*resp.Results) == 0 {
		t.Fatalf("expected results")
	}
	first := (*resp.Results)[0]
	if first.Ticker == nil || strings.EqualFold(*first.Ticker, "") {
		t.Fatalf("expected first result to have a ticker")
	}
}

func TestPolygonGetTickerDailyClose(t *testing.T) {
	if polygonConnection == nil {
		t.Skip("test server not initialized")
	}

	resp, err := polygonConnection.PolygonGetTickerDailyClose(testTicker)
	if err != nil {
		t.Fatalf("PolygonGetTickerDailyClose error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if resp.Results == nil || len(*resp.Results) != 1 {
		t.Fatalf("expected exactly one result for prev endpoint")
	}
	first := (*resp.Results)[0]
	if first.Close == nil && first.Open == nil {
		t.Fatalf("expected price fields to be present")
	}
	log.Println("Got response from PolygonGetTickerDailyClose:", PolygonResponseToString(resp))
}

func TestPolygonGetTickerHistory(t *testing.T) {
	if polygonConnection == nil {
		t.Skip("test server not initialized")
	}

	end := time.Now().UTC().AddDate(0, 0, -3)
	start := end.AddDate(0, 0, -7)
	resp, err := polygonConnection.PolygonGetTickerHistory(testTicker, start, end, 5)
	if err != nil {
		t.Fatalf("PolygonGetTickerHistory error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if resp.Results == nil || len(*resp.Results) == 0 {
		t.Fatalf("expected historical results")
	}
	for i, result := range *resp.Results {
		if !time.UnixMilli(*result.Timestamp).After(start) || !time.UnixMilli(*result.Timestamp).Before(end) {
			t.Fatalf("Result #%d in PolygonGetTickerHistory is from %s, which is not within range (%s - %s)", i, time.UnixMilli(*result.Timestamp).Format("2006-01-02T15:04:05Z"), start.Format("2006-01-02T15:04:05Z"), end.Format("2006-01-02T15:04:05Z"))
		}
	}
	log.Println("Got response from PolygonGetTickerHistory:", PolygonResponseToString(resp))
}

func TestPolygonGetTickerNews(t *testing.T) {
	if polygonConnection == nil {
		t.Skip("test server not initialized")
	}

	end := time.Now().UTC().AddDate(0, 0, -5)
	start := end.AddDate(0, 0, -10)
	resp, err := polygonConnection.PolygonGetTickerNews(testTicker, start, end, 10)
	if err != nil {
		t.Fatalf("PolygonGetTickerNews error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if resp.Results == nil || len(*resp.Results) == 0 {
		t.Fatalf("expected news results")
	}
	for i, result := range *resp.Results {
		if !result.PublishedUTC.After(start) || !result.PublishedUTC.Before(end) {
			t.Fatalf("Result #%d in PolygonGetTickerNews was published on %s, which is not within range (%s - %s)", i, result.PublishedUTC.Format("2006-01-02T15:04:05Z"), start.Format("2006-01-02T15:04:05Z"), end.Format("2006-01-02T15:04:05Z"))
		}
	}
	log.Println("Got response from PolygonGetTickerNews:", PolygonResponseToString(resp))
}

func TestPolygonGetTickerNews_SingleDay(t *testing.T) {
	if polygonConnection == nil {
		t.Skip("test server not initialized")
	}

	end := time.Now().UTC().AddDate(0, 0, -5)
	start := end.AddDate(0, 0, -1)
	resp, err := polygonConnection.PolygonGetTickerNews(testTicker, start, end, 10)
	if err != nil {
		t.Fatalf("PolygonGetTickerNews error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if resp.Results == nil || len(*resp.Results) == 0 {
		t.Fatalf("expected news results")
	}
	for i, result := range *resp.Results {
		if !result.PublishedUTC.After(start) || !result.PublishedUTC.Before(end) {
			t.Fatalf("Result #%d in PolygonGetTickerNews was published on %s, which is not within range (%s - %s)", i, result.PublishedUTC.Format("2006-01-02T15:04:05Z"), start.Format("2006-01-02T15:04:05Z"), end.Format("2006-01-02T15:04:05Z"))
		}
	}
	log.Println("Got response from PolygonGetTickerNews:", PolygonResponseToString(resp))
}

func TestPolygonGetTickerNews_SingleWeek(t *testing.T) {
	if polygonConnection == nil {
		t.Skip("test server not initialized")
	}

	end := time.Now().UTC().AddDate(0, 0, -5)
	start := end.AddDate(0, 0, -30)
	resp, err := polygonConnection.PolygonGetTickerNews(testTicker, start, end, 500)
	if err != nil {
		t.Fatalf("PolygonGetTickerNews error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if resp.Results == nil || len(*resp.Results) == 0 {
		t.Fatalf("expected news results")
	}
	for i, result := range *resp.Results {
		if !result.PublishedUTC.After(start) || !result.PublishedUTC.Before(end) {
			t.Fatalf("Result #%d in PolygonGetTickerNews was published on %s, which is not within range (%s - %s)", i, result.PublishedUTC.Format("2006-01-02T15:04:05Z"), start.Format("2006-01-02T15:04:05Z"), end.Format("2006-01-02T15:04:05Z"))
		}
	}
	log.Fatal("Got response from PolygonGetTickerNews:", PolygonResponseToString(resp))
}
