import React from 'react';
import './App.css';
import StockGraph from './components/stock-graph';
import SearchBar from './components/search-bar'

const InvestmentDashboard = () => {
  return (
    <div className="dashboard">
      {/* Main Content Area */}
      <div className="main-content">
        {/* Stock Graph Section */}
        <div className="graph-section">
          <h2>Portfolio Overview</h2>
          <StockGraph />
        </div>

        {/* Stock List Section */}
        <div className="stock-list-section">
          <h2>Watchlist</h2>
          <div className="stock-list">
            {['AAPL', 'GOOGL', 'MSFT', 'AMZN', 'META'].map((stock) => (
              <div key={stock} className="stock-item">
                <div className="stock-details">
                  <span className="stock-symbol">{stock}</span>
                  <span className="stock-price">$0.00</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      <SearchBar/>
    </div>
  );
};

export default InvestmentDashboard;