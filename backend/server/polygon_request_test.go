package server

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var testTicker string = "AMD"

var testServer *Server
var testTeardown func()
var testInitErr error

func TestMain(m *testing.M) {
	// Initialize server once for all tests
	testServer, testTeardown, testInitErr = func() (*Server, func(), error) {
		server, err := GetNewServer()
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

func SetupAndTeardown(tb *testing.TB, runServer bool) (*Server, func(), error) {
	server, err := GetNewServer()

	if err != nil {
		return nil, nil, err
	}

	if runServer {
		go server.Router.Run(":3333")
	}

	return server, func() {
		log.Println("Tearing down test environment...")
	}, err
}

func TestPolygonGetTicker(t *testing.T) {
	if testServer == nil {
		t.Skip("test server not initialized")
	}

	resp, err := testServer.PolygonGetTicker(testTicker)
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
	if testServer == nil {
		t.Skip("test server not initialized")
	}

	resp, err := testServer.PolygonGetTickerDailyClose(testTicker)
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
	if testServer == nil {
		t.Skip("test server not initialized")
	}

	end := time.Now().UTC().AddDate(0, 0, -3)
	start := end.AddDate(0, 0, -7)
	resp, err := testServer.PolygonGetTickerHistory(testTicker, start, end, 5)
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
	if testServer == nil {
		t.Skip("test server not initialized")
	}

	end := time.Now().UTC().AddDate(0, 0, -5)
	start := end.AddDate(0, 0, -10)
	resp, err := testServer.PolygonGetTickerNews(testTicker, start, end, 10)
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
	if testServer == nil {
		t.Skip("test server not initialized")
	}

	end := time.Now().UTC().AddDate(0, 0, -5)
	start := end.AddDate(0, 0, -1)
	resp, err := testServer.PolygonGetTickerNews(testTicker, start, end, 10)
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
