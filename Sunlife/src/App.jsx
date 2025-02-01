import React from 'react';
import './App.css';

const InvestmentDashboard = () => {
  return (
    <div className="dashboard">
      {/* Main Content Area */}
      <div className="main-content">
        {/* Stock Graph Section */}
        <div className="graph-section">
          <h2>Market Overview</h2>
          <div className="graph-placeholder">
            <span>Stock Graph Placeholder</span>
          </div>
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