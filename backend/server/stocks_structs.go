package server

// Returned by /api/v1/stocks/tickers/:symbol
type ServerTickerInfoResponse struct {
	Symbol          string  `json:"symbol"`
	Name            string  `json:"name"`
	Industry        string  `json:"industry"`
	Locale          string  `json:"locale"`
	PrimaryExchange string  `json:"primary_exchange"`
	OpenPrice       float64 `json:"open_price"`
	ClosePrice      float64 `json:"close_price"`
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

// Returned by /api/v1/stocks/tickers/:symbol/holdings
type TickerHoldings struct {
	Holdings []HoldingInfo `json:"holdings"`
}

type HoldingInfo struct {
	Symbol        string                   `json:"symbol"`
	CurrentShares float32                  `json:"current_shares"`
	History       []map[string]interface{} `json:"history"`
	ShareInfo     ServerTickerInfoResponse `json:"share_info"`
}

type Holding struct {
	Symbol        string  `json:"symbol"`
	CurrentShares float32 `json:"current_shares"`
}

// Dummy data
type StockTransaction struct {
	Symbol      string  `json:"symbol"`
	TotalShares float32 `json:"total_shares"`
	ShareChange float32 `json:"share_change"`
	Date        int64   `json:"date"`
}

var testTickerPurchases = []StockTransaction{
	{Symbol: "AAPL", TotalShares: 0.05, ShareChange: 0.05, Date: 1721985246},
	{Symbol: "AAPL", TotalShares: 0.1, ShareChange: 0.05, Date: 1722285246},
	{Symbol: "AAPL", TotalShares: 0.15, ShareChange: 0.05, Date: 1722585246},
	{Symbol: "AAPL", TotalShares: 0.1, ShareChange: -0.05, Date: 1728885246},
	{Symbol: "GOOGL", TotalShares: 0.1, ShareChange: 0.1, Date: 1723185246},
	{Symbol: "GOOGL", TotalShares: 0.2, ShareChange: 0.1, Date: 1723485246},
	{Symbol: "GOOGL", TotalShares: 0.15, ShareChange: -0.05, Date: 1723785246},
	{Symbol: "MSFT", TotalShares: 0.15, ShareChange: 0.15, Date: 1724085246},
	{Symbol: "MSFT", TotalShares: 0.3, ShareChange: 0.15, Date: 1724385246},
	{Symbol: "MSFT", TotalShares: 0.25, ShareChange: -0.05, Date: 1724685246},
	{Symbol: "AMZN", TotalShares: 0.2, ShareChange: 0.2, Date: 1733985246},
	{Symbol: "AMZN", TotalShares: 0.4, ShareChange: 0.2, Date: 1725285246},
	{Symbol: "AMZN", TotalShares: 0.35, ShareChange: -0.05, Date: 1733085246},
	{Symbol: "TSLA", TotalShares: 0.25, ShareChange: 0.25, Date: 1725885246},
	{Symbol: "TSLA", TotalShares: 0.5, ShareChange: 0.25, Date: 1726185246},
	{Symbol: "TSLA", TotalShares: 0.45, ShareChange: -0.05, Date: 1726485246},
	{Symbol: "META", TotalShares: 0.075, ShareChange: 0.075, Date: 1726785246},
	{Symbol: "META", TotalShares: 0.15, ShareChange: 0.075, Date: 1730685246},
	{Symbol: "META", TotalShares: 0.125, ShareChange: -0.025, Date: 1727385246},
	{Symbol: "NFLX", TotalShares: 0.125, ShareChange: 0.125, Date: 1727685246},
	{Symbol: "NFLX", TotalShares: 0.25, ShareChange: 0.125, Date: 1727985246},
	{Symbol: "NFLX", TotalShares: 0.225, ShareChange: -0.025, Date: 1728285246},
	{Symbol: "NVDA", TotalShares: 0.175, ShareChange: 0.175, Date: 1728585246},
	{Symbol: "NVDA", TotalShares: 0.35, ShareChange: 0.175, Date: 1728885246},
	{Symbol: "NVDA", TotalShares: 0.325, ShareChange: -0.025, Date: 1729185246},
	{Symbol: "BABA", TotalShares: 0.225, ShareChange: 0.225, Date: 1729485246},
	{Symbol: "BABA", TotalShares: 0.45, ShareChange: 0.225, Date: 1729785246},
	{Symbol: "BABA", TotalShares: 0.425, ShareChange: -0.025, Date: 1730085246},
	{Symbol: "V", TotalShares: 0.275, ShareChange: 0.275, Date: 1730385246},
	{Symbol: "V", TotalShares: 0.55, ShareChange: 0.275, Date: 1730685246},
	{Symbol: "V", TotalShares: 0.525, ShareChange: -0.025, Date: 1722585246},
	{Symbol: "JPM", TotalShares: 0.05, ShareChange: 0.05, Date: 1731285246},
	{Symbol: "JPM", TotalShares: 0.1, ShareChange: 0.05, Date: 1731585246},
	{Symbol: "JPM", TotalShares: 0.075, ShareChange: -0.025, Date: 1731885246},
	{Symbol: "JNJ", TotalShares: 0.1, ShareChange: 0.1, Date: 1732185246},
	{Symbol: "JNJ", TotalShares: 0.2, ShareChange: 0.1, Date: 1732485246},
	{Symbol: "JNJ", TotalShares: 0.175, ShareChange: -0.025, Date: 1732785246},
	{Symbol: "WMT", TotalShares: 0.15, ShareChange: 0.15, Date: 1733085246},
	{Symbol: "WMT", TotalShares: 0.3, ShareChange: 0.15, Date: 1733385246},
	{Symbol: "WMT", TotalShares: 0.275, ShareChange: -0.025, Date: 1733685246},
	{Symbol: "PG", TotalShares: 0.2, ShareChange: 0.2, Date: 1733985246},
	{Symbol: "PG", TotalShares: 0.4, ShareChange: 0.2, Date: 1734285246},
	{Symbol: "PG", TotalShares: 0.375, ShareChange: -0.025, Date: 1734585246},
	{Symbol: "DIS", TotalShares: 0.25, ShareChange: 0.25, Date: 1734885246},
	{Symbol: "DIS", TotalShares: 0.5, ShareChange: 0.25, Date: 1735185246},
	{Symbol: "DIS", TotalShares: 0.475, ShareChange: -0.025, Date: 1735485246},
	{Symbol: "MA", TotalShares: 0.075, ShareChange: 0.075, Date: 1735785246},
	{Symbol: "MA", TotalShares: 0.15, ShareChange: 0.075, Date: 1736085246},
	{Symbol: "MA", TotalShares: 0.125, ShareChange: -0.025, Date: 1736385246},
	{Symbol: "HD", TotalShares: 0.125, ShareChange: 0.125, Date: 1736685246},
	{Symbol: "HD", TotalShares: 0.25, ShareChange: 0.125, Date: 1736985246},
	{Symbol: "HD", TotalShares: 0.225, ShareChange: -0.025, Date: 1737285246},
	{Symbol: "VZ", TotalShares: 0.175, ShareChange: 0.175, Date: 1735478270},
	{Symbol: "VZ", TotalShares: 0.35, ShareChange: 0.175, Date: 1737885246},
	{Symbol: "VZ", TotalShares: 0.325, ShareChange: -0.025, Date: 1737585246},
	{Symbol: "PYPL", TotalShares: 0.225, ShareChange: 0.225, Date: 1733085246},
	{Symbol: "PYPL", TotalShares: 0.45, ShareChange: 0.225, Date: 1729485246},
	{Symbol: "PYPL", TotalShares: 0.425, ShareChange: -0.025, Date: 1730685246},
	{Symbol: "ADBE", TotalShares: 0.275, ShareChange: 0.275, Date: 1734885246},
	{Symbol: "ADBE", TotalShares: 0.55, ShareChange: 0.275, Date: 1729185246},
	{Symbol: "ADBE", TotalShares: 0.525, ShareChange: -0.025, Date: 1734285246},
	{Symbol: "INTC", TotalShares: 0.3, ShareChange: 0.3, Date: 1733585246},
	{Symbol: "INTC", TotalShares: 0.6, ShareChange: 0.3, Date: 1736478270},
	{Symbol: "INTC", TotalShares: 0.575, ShareChange: -0.025, Date: 1733185246},
	{Symbol: "CSCO", TotalShares: 0.35, ShareChange: 0.35, Date: 1735485246},
	{Symbol: "CSCO", TotalShares: 0.7, ShareChange: 0.35, Date: 1734285246},
	{Symbol: "CSCO", TotalShares: 0.675, ShareChange: -0.025, Date: 1730985246},
	{Symbol: "ORCL", TotalShares: 0.4, ShareChange: 0.4, Date: 1732478270},
	{Symbol: "ORCL", TotalShares: 0.8, ShareChange: 0.4, Date: 1724885246},
	{Symbol: "ORCL", TotalShares: 0.775, ShareChange: -0.025, Date: 1726485246},
}

/* var testTickerHistory = TickerHistory{
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
*/
