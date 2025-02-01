import React from 'react';
import './App.css';
import StockGraph from './components/stock-graph';
import TickerList from './components/ticker-list';

const InvestmentDashboard = () => {
  return (
    <div className="dashboard">
      {/* Main Content Area */}
      <div className="main-content">
        {/* Stock Graph Section */}
        <div className="graph-section">
          <h2>Market Overview</h2>
          <StockGraph />
        </div>

        {/* Stock List Section */}
        <div className="stock-list-section">
          <h2>Watchlist</h2>
          <div className="stock-list">
            <TickerList />
          </div>
        </div>
      </div>

      {/* Right Sidebar */}
      <div className="sidebar">
        <div className="search-container">
          <input
            type="text"
            placeholder="Search stocks..."
            className="search-input"
          />
          <svg
            className="search-icon"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
      </div>
    </div>
  );
};

export default InvestmentDashboard;