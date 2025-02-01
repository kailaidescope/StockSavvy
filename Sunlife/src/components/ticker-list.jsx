import React, { useState } from 'react'
import { BsFillChatRightTextFill } from "react-icons/bs";
import TickerGraph from './ticker-graph';
import { useSymbol } from '../contexts/symbol-context';

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
    trend: "+0.01",
    emoji: '🥶'
  },
  {
    symbol: "MSFT",
    price: "10.00",
    trend: "+0.01",
    emoji: '🔥'
  },
]

const TickerList = () => {
  const [expandedStock, setExpandedStock] = useState(null);
  const {selectedSymbol, setSelectedSymbol} = useSymbol();

  const handleClickItem = (symbol) => {
    setExpandedStock(expandedStock === symbol ? null : symbol);
  }

  const handleClickChat = (event, symbol) => {
    event.stopPropagation();
    setSelectedSymbol(symbol);
  }

  return (
    <div>
      {mockData.map((stock) => (
        <div key={stock.symbol} className="stock-item" onClick={() => handleClickItem(stock.symbol)}>
          <div className="stock-details">
            <span className="stock-symbol">{stock?.emoji} {stock.symbol}</span>
            <div className='stock-price-details'>
              <span className="stock-price">${stock.price}</span>
              <span className="stock-trend">{stock.trend}</span>
              <BsFillChatRightTextFill className='chat' color='grey' onClick={(event) => handleClickChat(event, stock.symbol)}/>
            </div>
          </div>
          {expandedStock === stock.symbol && (
            <div className="stock-graph">
              <TickerGraph />
            </div>
          )}
        </div>
      ))}
      <style jsx>
        {`
        .stock-item {
    padding: 12px;
    border-bottom: 1px solid grey;
    cursor: pointer;
}

.stock-item:hover {
    background-color: #fdf0c4
}

.stock-item:active {
    background-color: var(--color-jonquil)
}

.stock-details {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.stock-emoji {
    flex: 1;
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

chat {
    z-index:1000
}

.stock-graph {
    padding: 12px;
    background-color: #f9f9f9;
    border-top: 1px solid grey;
}
`}
      </style>
    </div>
  )
}

export default TickerList