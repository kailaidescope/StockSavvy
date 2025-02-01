import React from 'react';
import '../App.css';
import StockGraph from './stock-graph';
import TickerList from './ticker-list';
import SearchBar from './search-bar'

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
            <TickerList />
          </div>
        </div>
      </div>

      <SearchBar/>
    </div>
  );
};

export default InvestmentDashboard;