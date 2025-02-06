package server

// This file contains the functions that make requests to the Polygon API
// https://polygon.io/docs/stocks/getting-started
//
// All calls to Polygon should be directed through these functions.
// If a new call is needed, it shoud be added to this file.

// For making JSONs
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GenericPolygonRequest[T any](url string) (*T, error) {
	time.Sleep(THROTTLE_TIME * time.Second)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, errors.Join(errors.New("error generating request for polygon.io"), err)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("error sending/receiving request to polygon.io"), err)
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("error: %d", res.StatusCode))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("error reading polygon.io response"), err)
	}

	// Unmarshall the unmarshalledBody
	var unmarshalledBody *T
	if err = json.Unmarshal(body, unmarshalledBody); err != nil {
		return nil, errors.Join(errors.New("error unmarshalling response"), err)
	}

	return unmarshalledBody, nil
}

type PolygonGetTickerResponse struct {
	Results []struct {
		Ticker          *string   `json:"ticker"`
		Name            *string   `json:"name"`
		Market          *string   `json:"market"`
		Locale          *string   `json:"locale"`
		PrimaryExchange *string   `json:"primary_exchange"`
		Type            *string   `json:"type"`
		Active          *bool     `json:"active"`
		CurrencyName    *string   `json:"currency_name"`
		Cik             *string   `json:"cik"`
		CompositeFigi   *string   `json:"composite_figi"`
		ShareClassFigi  *string   `json:"share_class_figi"`
		LastUpdatedUtc  time.Time `json:"last_updated_utc"`
	} `json:"results"`
	Status    *string `json:"status"`
	RequestID *string `json:"request_id"`
	Count     *int    `json:"count"`
}

// PolygonGetTicker returns the ticker information for a given symbol
//
// Input:
//   - symbol: the symbol of the stock
//
// Output:
//   - *GetTickerResponse: the response from the Polygon API
//   - error: any error that occurred
func (server *Server) PolygonGetTicker(symbol string) (*PolygonGetTickerResponse, error) {
	url := fmt.Sprintf("https://api.polygon.io/v3/reference/tickers?ticker=%s&active=true&limit=100&apiKey=%s", symbol, server.GetPolygonKey())

	response, err := GenericPolygonRequest[PolygonGetTickerResponse](url)
	if err != nil {
		return nil, errors.Join(errors.New("error getting info from polygon"), err)
	}
	if response.Results == nil || len(response.Results) == 0 || response.Count == nil || *response.Count == 0 {
		return nil, errors.New("no results found")
	}
	return response, nil
}

type PolygonGetTickerAggregateResponse struct {
	Ticker       string `json:"ticker"`
	QueryCount   int    `json:"queryCount"`
	ResultsCount int    `json:"resultsCount"`
	Adjusted     bool   `json:"adjusted"`
	Results      []struct {
		V  int     `json:"v"`
		Vw float64 `json:"vw"`
		O  float64 `json:"o"`
		C  float64 `json:"c"`
		H  float64 `json:"h"`
		L  float64 `json:"l"`
		T  int64   `json:"t"`
		N  int     `json:"n"`
	} `json:"results"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Count     int    `json:"count"`
	NextURL   string `json:"next_url"`
}

// PolygonGetTickerAggregate returns the ticker aggregate information for a given symbol for the last day
//
// Input:
//   - symbol: the symbol of the stock
//
// Output:
//   - *GetTickerAggregateResponse: the response from the Polygon API
//   - error: any error that occurred
func (server *Server) PolygonGetTickerLastHistory(symbol string) (*PolygonGetTickerAggregateResponse, error) {
	url := fmt.Sprintf("https://api.polygon.io/v3/reference/tickers?ticker=%s&active=true&limit=100&apiKey=%s", symbol, server.GetPolygonKey())
	return GenericPolygonRequest[PolygonGetTickerAggregateResponse](url)
}
