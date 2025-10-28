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
	"log"
	"net/http"
	"reflect"
	"time"
)

// Sends a customizable GET request to Polygon.io's API
//
// Input:
//   - url: the url to send the request to
//   - T: the type of the response
//
// Output:
//   - *T: the response from the Polygon API
//   - error: non-nil if an error occurred during the request or if the response was not 200
func GenericPolygonGetRequest[T any](url string) (*T, error) {
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
	// Check if the response is valid (status code 2XX)
	if !(res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices) {
		return nil, fmt.Errorf("error: %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Unmarshall the decodedBody
	var decodedBody T
	if err = json.NewDecoder(res.Body).Decode(&decodedBody); err != nil {
		return nil, errors.Join(errors.New("error decoding response"), err)
	}

	// TODO: Check if this function is necessary, or if it can be done in a better way. Maybe a less brute force way?
	isNil, err := anyFieldIsNil(&decodedBody)
	if err != nil {
		return nil, errors.Join(errors.New("error checking if decoded response has nil fields"), err)
	}
	if isNil {
		log.Println("Response has some nil fields")
		//return nil, errors.New("decoded response has nil fields")
	}

	return &decodedBody, nil
}

// Checks if any field in a struct is nil
func anyFieldIsNil(v interface{}) (bool, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else {
		return false, errors.New("value is not a pointer")
	}
	var nilFields []string
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Kind() == reflect.Ptr && val.Field(i).IsNil() {
			nilFields = append(nilFields, val.Type().Field(i).Name)
		}
	}
	if len(nilFields) > 0 {
		log.Printf("Nil fields: %v\n", nilFields)
		return true, nil
	}
	return false, nil
}

type PolygonGetTickerResponse struct {
	Results *[]struct {
		Ticker          *string    `json:"ticker"`
		Name            *string    `json:"name"`
		Market          *string    `json:"market"`
		Locale          *string    `json:"locale"`
		PrimaryExchange *string    `json:"primary_exchange"`
		Type            *string    `json:"type"`
		Active          *bool      `json:"active"`
		CurrencyName    *string    `json:"currency_name"`
		Cik             *string    `json:"cik"`
		CompositeFigi   *string    `json:"composite_figi"`
		ShareClassFigi  *string    `json:"share_class_figi"`
		LastUpdatedUtc  *time.Time `json:"last_updated_utc"`
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

	response, err := GenericPolygonGetRequest[PolygonGetTickerResponse](url)
	if err != nil {
		return nil, errors.Join(errors.New("error getting info from polygon"), err)
	}
	if response.Results == nil || len(*response.Results) == 0 || response.Count == nil || *response.Count == 0 {
		return nil, errors.New("no results found")
	}

	return response, nil
}

type PolygonGetTickerAggregateResponse struct {
	Ticker       *string `json:"ticker"`
	QueryCount   *int    `json:"queryCount"`
	ResultsCount *int    `json:"resultsCount"`
	Adjusted     *bool   `json:"adjusted"`
	Results      *[]struct {
		Ticker       *string  `json:"T"`
		Volume       *float64 `json:"v"`
		VWAP         *float64 `json:"vw"`
		Open         *float64 `json:"o"`
		Close        *float64 `json:"c"`
		High         *float64 `json:"h"`
		Low          *float64 `json:"l"`
		Timestamp    *int64   `json:"t"`
		Transactions *int     `json:"n"`
		OTC          *bool    `json:"otc"` // This is omitted if false
	} `json:"results"`
	Status    *string `json:"status"`
	RequestID *string `json:"request_id"`
	Count     *int    `json:"count"`
}

// PolygonGetTickerAggregate returns the ticker aggregate information for a given symbol for the last day
//
// Input:
//   - symbol: the symbol of the stock
//
// Output:
//   - *GetTickerAggregateResponse: the response from the Polygon API
//   - error: any error that occurred
func (server *Server) PolygonGetTickerDailyClose(symbol string) (*PolygonGetTickerAggregateResponse, error) {
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/prev?apiKey=%s", symbol, server.GetPolygonKey())
	response, err := GenericPolygonGetRequest[PolygonGetTickerAggregateResponse](url)
	if err != nil {
		return nil, errors.Join(errors.New("error getting info from polygon"), err)
	}
	if response.Results == nil || len(*response.Results) == 0 || response.Count == nil || *response.Count == 0 {
		return nil, errors.New("no results found")
	}
	if len(*response.Results) != 1 {
		return nil, errors.New("too many results found")
	}

	return response, nil
}

type PolygonGetTickerHistoryResponse struct {
	Ticker       *string `json:"ticker"`
	QueryCount   *int    `json:"queryCount"`
	ResultsCount *int    `json:"resultsCount"`
	Adjusted     *bool   `json:"adjusted"`
	Results      *[]struct {
		Volume       *float64 `json:"v"`
		VWAP         *float64 `json:"vw"`
		Open         *float64 `json:"o"`
		Close        *float64 `json:"c"`
		High         *float64 `json:"h"`
		Low          *float64 `json:"l"`
		Timestamp    *int64   `json:"t"`
		Transactions *int     `json:"n"`
		OTC          *bool    `json:"otc"`
	} `json:"results"`
	Status    *string `json:"status"`
	RequestID *string `json:"request_id"`
	Count     *int    `json:"count"`
}

// PolygonGetTickerHistory returns the ticker's historical price data within the selected time range
//
// Input:
//   - symbol: the symbol of the stock
//   - startDate: the earliest date to retrieve data from
//   - endDate: the lastest date to retrieve data from
//   - limit: the limit of days to retrieve, set <= 0 for maximum
//
// Output:
//   - *GetTickerAggregateResponse: the response from the Polygon API
//   - error: any error that occurred
func (server *Server) PolygonGetTickerHistory(symbol string, startDate time.Time, endDate time.Time, limit int) (*PolygonGetTickerHistoryResponse, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Set default limit
	responseLengthLimit := 5000
	if limit > 0 {
		responseLengthLimit = limit
	}
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/1/day/%s/%s?adjusted=true&sort=asc&limit=%d&apiKey=%s", symbol, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), responseLengthLimit, server.GetPolygonKey())
	response, err := GenericPolygonGetRequest[PolygonGetTickerHistoryResponse](url)
	if err != nil {
		return nil, errors.Join(errors.New("error getting info from polygon"), err)
	}
	if response.Results == nil || len(*response.Results) == 0 || response.Count == nil || *response.Count == 0 {
		return nil, errors.New("no results found")
	}

	return response, nil
}

type PolygonGetTickerNews struct {
	Results *[]struct {
		ID        *string `json:"id"`
		Publisher *struct {
			Name        *string `json:"name"`
			HomepageURL *string `json:"homepage_url"`
			LogoURL     *string `json:"logo_url"`
			FaviconURL  *string `json:"favicon_url"`
		} `json:"publisher"`
		Title        *string    `json:"title"`
		Author       *string    `json:"author"`
		PublishedUTC *time.Time `json:"published_utc"`
		ArticleURL   *string    `json:"article_url"`
		Tickers      *[]string  `json:"tickers"`
		ImageURL     *string    `json:"image_url"`
		Description  *string    `json:"description"`
		Keywords     *[]string  `json:"keywords"`
		Insights     *[]struct {
			Ticker             *string `json:"ticker"`
			Sentiment          *string `json:"sentiment"`
			SentimentReasoning *string `json:"sentiment_reasoning"`
		} `json:"insights"`
	} `json:"results"`
	Status    *string `json:"status"`
	RequestID *string `json:"request_id"`
	Count     *int    `json:"count"`
	NextURL   *string `json:"next_url"`
}

// PolygonGetTickerNews returns the ticker's historical price data within the selected time range
//
// Input:
//   - symbol: the symbol of the stock
//   - startDate: the earliest date to retrieve data from
//   - endDate: the lastest date to retrieve data from
//   - limit: the limit of days to retrieve, set <= 0 for maximum
//
// Output:
//   - *GetTickerAggregateResponse: the response from the Polygon API
//   - error: any error that occurred
func (server *Server) PolygonGetTickerNews(symbol string, startDate time.Time, endDate time.Time, limit int) (*PolygonGetTickerNews, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Set default limit
	responseLengthLimit := 300
	if limit > 0 {
		responseLengthLimit = limit
	}
	order := "desc"
	sort := "published_utc"
	url := fmt.Sprintf("https://api.polygon.io/v2/reference/news?ticker=%s&order=%s&limit=%d&sort=%s&apiKey=%s&published_utc.gte=%s&published_utc.lte=%s", symbol, order, responseLengthLimit, sort, server.GetPolygonKey(), startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z"))
	response, err := GenericPolygonGetRequest[PolygonGetTickerNews](url)
	if err != nil {
		return nil, errors.Join(errors.New("error getting info from polygon"), err)
	}
	if response.Results == nil || len(*response.Results) == 0 || response.Count == nil || *response.Count == 0 {
		return nil, errors.New("no results found")
	}

	return response, nil
}
