import React from 'react'
import './ticker-list.css'

const mockData = [
  {
    symbol: "AAPL",
    price: "10.00",
    trend: "+0.01"
  },
  {
    symbol: "GOOGL",
    price: "10.00",
    trend: "+0.01"
  },
  {
    symbol: "AMZN",
    price: "10.00",
    trend: "+0.01"
  },
  {
    symbol: "MSFT",
    price: "10.00",
    trend: "+0.01"
  },
]

const TickerList = () => {
  return (
    <div>
      {mockData.map((stock) => (
        <div key={stock} className="stock-item">
          <div className="stock-details">
            <span className="stock-symbol">{stock.symbol}</span>
            <div className='stock-price-details'>
              <span className="stock-price">${stock.price}</span>
              <span className="stock-trend">{stock.trend}</span>
            </div>
          </div>
        </div>
      ))}
      <style jsx>
        {`
        .stock-item {
    padding: 12px;
    border-bottom: 1px solid #f3f4f6;
    cursor: pointer;
}

.stock-item:hover {
    background-color: #f9fafb;
}

.stock-details {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.stock-symbol {
    flex: 1;
    font-weight: 500;
    color: #111827;
}

.stock-price {
    color: #6b7280;
}

.stock-trend {
    color: #6b7280;
}

.stock-price-details {
    flex: 1;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

/* Scrollbar Styling */
.stock-list::-webkit-scrollbar {
    width: 6px;
}

.stock-list::-webkit-scrollbar-track {
    background: #f1f1f1;
}

.stock-list::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 3px;
}

.stock-list::-webkit-scrollbar-thumb:hover {
    background: #555;
}
`}
      </style>
    </div>
  )
}

export default TickerList