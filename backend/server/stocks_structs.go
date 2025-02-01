package server

// Returned by /api/v1/stocks/tickers/:symbol
type TickerInfo struct {
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	Industry        string `json:"industry"`
	Locale          string `json:"locale"`
	PrimaryExchange string `json:"primary_exchange"`
}

// Returned by /api/v1/stocks/tickers/:symbol/history
type TickerHistory struct {
	History []map[string]interface{} `json:"history"`
}

// Returned by /api/v1/stocks/tickers/:symbol
type TickerNews struct {
	AverageSentiment float32 `json:"avg_sentiment"`
	StdDevSentiment  float32 `json:"std_dev_sentiment"`
	NumArticles      int     `json:"num_articles"`
}

// Dummy data
var testTickerHistory = TickerHistory{
	History: []map[string]interface{}{
		{"time": 1612137600, "value": 100.0},
		{"time": 1612224000, "value": 101.5},
		{"time": 1612310400, "value": 102.3},
		{"time": 1612396800, "value": 103.8},
		{"time": 1612483200, "value": 104.2},
		{"time": 1612569600, "value": 105.0},
		{"time": 1612656000, "value": 106.1},
		{"time": 1612742400, "value": 107.3},
		{"time": 1612828800, "value": 108.5},
		{"time": 1612915200, "value": 109.7},
		{"time": 1613001600, "value": 110.2},
		{"time": 1613088000, "value": 111.4},
		{"time": 1613174400, "value": 112.6},
		{"time": 1613260800, "value": 113.8},
		{"time": 1613347200, "value": 114.0},
		{"time": 1613433600, "value": 115.2},
		{"time": 1613520000, "value": 116.4},
		{"time": 1613606400, "value": 117.6},
		{"time": 1613692800, "value": 118.8},
		{"time": 1613779200, "value": 119.0},
		{"time": 1613865600, "value": 120.2},
		{"time": 1613952000, "value": 121.4},
		{"time": 1614038400, "value": 122.6},
		{"time": 1614124800, "value": 123.8},
		{"time": 1614211200, "value": 124.0},
		{"time": 1614297600, "value": 125.2},
		{"time": 1614384000, "value": 126.4},
		{"time": 1614470400, "value": 127.6},
		{"time": 1614556800, "value": 128.8},
		{"time": 1614643200, "value": 129.0},
		{"time": 1614729600, "value": 130.2},
		{"time": 1614816000, "value": 131.4},
		{"time": 1614902400, "value": 132.6},
		{"time": 1614988800, "value": 133.8},
		{"time": 1615075200, "value": 134.0},
		{"time": 1615161600, "value": 135.2},
		{"time": 1615248000, "value": 136.4},
		{"time": 1615334400, "value": 137.6},
		{"time": 1615420800, "value": 138.8},
		{"time": 1615507200, "value": 139.0},
		{"time": 1615593600, "value": 140.2},
		{"time": 1615680000, "value": 141.4},
		{"time": 1615766400, "value": 142.6},
		{"time": 1615852800, "value": 143.8},
		{"time": 1615939200, "value": 144.0},
		{"time": 1616025600, "value": 145.2},
		{"time": 1616112000, "value": 146.4},
		{"time": 1616198400, "value": 147.6},
		{"time": 1616284800, "value": 148.8},
		{"time": 1616371200, "value": 149.0},
	},
}
