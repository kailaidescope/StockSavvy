package polygon

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
	"os"
	"reflect"
	"strings"
	"time"
)

var errLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

type PolygonConnection struct {
	polygonKeys       []string
	currentPolygonKey int
	ThrottleTime      time.Duration
}

func GetPolygonConnection(polygonKeys []string, throttleTime time.Duration) *PolygonConnection {
	return &PolygonConnection{
		polygonKeys:       polygonKeys,
		currentPolygonKey: 0,
		ThrottleTime:      throttleTime,
	}
}

func (polygonConnection *PolygonConnection) GetPolygonKey() string {
	polygonConnection.currentPolygonKey = (polygonConnection.currentPolygonKey + 1) % len(polygonConnection.polygonKeys)
	return polygonConnection.polygonKeys[polygonConnection.currentPolygonKey]
}

// Sends a customizable GET request to Polygon.io's API
//
// Input:
//   - url: the url to send the request to
//   - T: the type of the response
//
// Output:
//   - *T: the response from the Polygon API
//   - error: non-nil if an error occurred during the request or if the response was not 200
func GenericPolygonGetRequest[T any](polygonConnection *PolygonConnection, url string) (*T, error) {
	time.Sleep(polygonConnection.ThrottleTime)
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
		errLogger.Println("Response has some nil fields")
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
func (polygonConnection *PolygonConnection) PolygonGetTicker(symbol string) (*PolygonGetTickerResponse, error) {
	url := fmt.Sprintf("https://api.polygon.io/v3/reference/tickers?ticker=%s&active=true&limit=100&apiKey=%s", symbol, polygonConnection.GetPolygonKey())

	response, err := GenericPolygonGetRequest[PolygonGetTickerResponse](polygonConnection, url)
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
func (polygonConnection *PolygonConnection) PolygonGetTickerDailyClose(symbol string) (*PolygonGetTickerAggregateResponse, error) {
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/prev?apiKey=%s", symbol, polygonConnection.GetPolygonKey())
	response, err := GenericPolygonGetRequest[PolygonGetTickerAggregateResponse](polygonConnection, url)
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
func (polygonConnection *PolygonConnection) PolygonGetTickerHistory(symbol string, startDate time.Time, endDate time.Time, limit int) (*PolygonGetTickerHistoryResponse, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Set default limit
	responseLengthLimit := 5000
	if limit > 0 {
		responseLengthLimit = limit
	}
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/1/day/%s/%s?adjusted=true&sort=asc&limit=%d&apiKey=%s", symbol, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), responseLengthLimit, polygonConnection.GetPolygonKey())
	response, err := GenericPolygonGetRequest[PolygonGetTickerHistoryResponse](polygonConnection, url)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error getting info from polygon (response: %v)", response), err)
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
func (polygonConnection *PolygonConnection) PolygonGetTickerNews(symbol string, startDate time.Time, endDate time.Time, limit int) (*PolygonGetTickerNews, error) {
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
	url := fmt.Sprintf("https://api.polygon.io/v2/reference/news?ticker=%s&order=%s&limit=%d&sort=%s&apiKey=%s&published_utc.gte=%s&published_utc.lte=%s", symbol, order, responseLengthLimit, sort, polygonConnection.GetPolygonKey(), startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z"))
	response, err := GenericPolygonGetRequest[PolygonGetTickerNews](polygonConnection, url)
	if err != nil {
		return nil, errors.Join(errors.New("error getting info from polygon"), err)
	}
	if response.Results == nil || len(*response.Results) == 0 || response.Count == nil || *response.Count == 0 {
		return nil, errors.New("no results found")
	}

	return response, nil
}

// PolygonResponseToString converts any of the Polygon response structs in this file into a readable string using reflection.
// It attempts to pretty-print pointer fields, slices, structs and time.Time values.
func PolygonResponseToString(v interface{}) string {
	if v == nil {
		return ""
	}
	var b strings.Builder

	var formatValue func(reflect.Value, int)
	indentStr := func(level int) string { return strings.Repeat("  ", level) }

	formatValue = func(rv reflect.Value, indent int) {
		if !rv.IsValid() {
			b.WriteString(indentStr(indent) + "<invalid>\n")
			return
		}
		// Deref pointers
		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				b.WriteString(indentStr(indent) + "<nil>\n")
				return
			}
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Struct:
			typ := rv.Type()
			// Special-case time.Time
			if typ.PkgPath() == "time" && typ.Name() == "Time" {
				t := rv.Interface().(time.Time)
				b.WriteString(indentStr(indent) + t.Format(time.RFC3339) + "\n")
				return
			}
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldType := typ.Field(i)
				name := fieldType.Name
				// Print field name
				b.WriteString(indentStr(indent) + name + ": ")
				if !field.IsValid() {
					b.WriteString("<invalid>\n")
					continue
				}
				// If pointer -> handle inside
				if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						b.WriteString("<nil>\n")
						continue
					}
					elem := field.Elem()
					// time.Time pointer
					if elem.Kind() == reflect.Struct && elem.Type().PkgPath() == "time" && elem.Type().Name() == "Time" {
						t := elem.Interface().(time.Time)
						b.WriteString(t.Format(time.RFC3339) + "\n")
						continue
					}
					switch elem.Kind() {
					case reflect.String:
						b.WriteString(fmt.Sprintf("%s\n", elem.String()))
					case reflect.Bool:
						b.WriteString(fmt.Sprintf("%t\n", elem.Bool()))
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						b.WriteString(fmt.Sprintf("%d\n", elem.Int()))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						b.WriteString(fmt.Sprintf("%d\n", elem.Uint()))
					case reflect.Float32, reflect.Float64:
						b.WriteString(fmt.Sprintf("%v\n", elem.Float()))
					case reflect.Slice, reflect.Array:
						if elem.Len() == 0 {
							b.WriteString("[]\n")
							continue
						}
						b.WriteString("\n")
						for j := 0; j < elem.Len(); j++ {
							b.WriteString(indentStr(indent+1) + "- ")
							// print element (may be basic or struct)
							// capture current length to check if nested added newline
							startLen := b.Len()
							formatValue(elem.Index(j), indent+2)
							// if nested printed with its own newline, avoid double newlines; ensure at least one newline
							if b.Len() == startLen {
								b.WriteString("\n")
							}
						}
					case reflect.Struct:
						b.WriteString("\n")
						formatValue(elem, indent+1)
					default:
						// Fallback to fmt
						b.WriteString(fmt.Sprintf("%v\n", elem.Interface()))
					}
				} else {
					// Non-pointer field kinds
					switch field.Kind() {
					case reflect.Struct:
						// time.Time?
						if field.Type().PkgPath() == "time" && field.Type().Name() == "Time" {
							t := field.Interface().(time.Time)
							b.WriteString(t.Format(time.RFC3339) + "\n")
						} else {
							b.WriteString("\n")
							formatValue(field, indent+1)
						}
					case reflect.Slice, reflect.Array:
						if field.Len() == 0 {
							b.WriteString("[]\n")
							continue
						}
						b.WriteString("\n")
						for j := 0; j < field.Len(); j++ {
							b.WriteString(indentStr(indent+1) + "- ")
							startLen := b.Len()
							formatValue(field.Index(j), indent+2)
							if b.Len() == startLen {
								b.WriteString("\n")
							}
						}
					case reflect.String:
						b.WriteString(fmt.Sprintf("%s\n", field.String()))
					case reflect.Bool:
						b.WriteString(fmt.Sprintf("%t\n", field.Bool()))
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						b.WriteString(fmt.Sprintf("%d\n", field.Int()))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						b.WriteString(fmt.Sprintf("%d\n", field.Uint()))
					case reflect.Float32, reflect.Float64:
						b.WriteString(fmt.Sprintf("%v\n", field.Float()))
					case reflect.Interface:
						if field.IsNil() {
							b.WriteString("<nil>\n")
						} else {
							b.WriteString("\n")
							formatValue(field.Elem(), indent+1)
						}
					default:
						b.WriteString(fmt.Sprintf("%v\n", field.Interface()))
					}
				}
			}
		case reflect.Slice, reflect.Array:
			if rv.Len() == 0 {
				b.WriteString(indentStr(indent) + "[]\n")
				return
			}
			for i := 0; i < rv.Len(); i++ {
				b.WriteString(indentStr(indent) + "- ")
				startLen := b.Len()
				formatValue(rv.Index(i), indent+1)
				if b.Len() == startLen {
					b.WriteString("\n")
				}
			}
		default:
			// Basic types
			b.WriteString(indentStr(indent) + fmt.Sprintf("%v\n", rv.Interface()))
		}
	}

	formatValue(reflect.ValueOf(v), 0)
	return b.String()
}
