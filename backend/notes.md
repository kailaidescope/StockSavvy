# Notes for backend implementation

### Stock API

- [Polygon.io](https://polygon.io/docs/stocks/getting-started)
  - [Get ticker](https://polygon.io/docs/stocks/get_v3_reference_tickers)
  - [Get news about a ticker](https://polygon.io/docs/stocks/get_v2_reference_news)
  - [Daily open/close](https://polygon.io/docs/stocks/get_v1_open-close__stocksticker___date)
  - [Simple Moving avg](https://polygon.io/docs/stocks/get_v1_indicators_sma__stockticker)
  - [Exp moving avg](https://polygon.io/docs/stocks/get_v1_indicators_ema__stockticker)
- [NYT](https://developer.nytimes.com/apis)
  - Rate limit:
    > Yes, there are two rate limits per API: 500 requests per day and 5 requests per minute. You should sleep 12 seconds between calls to avoid hitting the per minute rate limit. If you need a higher rate limit, please contact us at code@nytimes.com.
  - [Articles by month and year](https://developer.nytimes.com/docs/archive-product/1/routes/%7Byear%7D/%7Bmonth%7D.json/get)
