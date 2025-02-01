import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import DropdownMenu from './dropdown-menu';
import TextBox from './text-box';

const fakeStocks = [
  'Apple (AAPL)',
  'Google (GOOGL)',
  'Microsoft (MSFT)',
  'Amazon (AMZN)',
  'Meta (META)',
  'Tesla (TSLA)',
  'Netflix (NFLX)',
  'Nvidia (NVDA)',
  'Adobe (ADBE)',
  'Intel (INTC)',
];

const SearchBar = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [filteredStocks, setFilteredStocks] = useState([]);
  const navigate = useNavigate();

  const handleSearch = (event) => {
    const value = event.target.value;
    setSearchTerm(value);
    if (value) {
      const filtered = fakeStocks.filter((stock) =>
        stock.toLowerCase().includes(value.toLowerCase())
      );
      setFilteredStocks(filtered);
    } else {
      setFilteredStocks([]);
    }
  };

  const handleStockClick = (stock) => {
    navigate(`/search/${encodeURIComponent(stock)}`);
  };

  return (
    <>
      <div className="sidebar">
        <div className="search-container">
          <input
            type="text"
            placeholder="Search stocks..."
            className="search-input"
            value={searchTerm}
            onChange={handleSearch}
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
          <DropdownMenu />
        </div>
        {filteredStocks.length > 0 && (
          <ul className="search-results">
            {filteredStocks.map((stock, index) => (
              <li
                key={index}
                className="search-result-item"
                onClick={() => handleStockClick(stock)}
              >
                {stock}
              </li>
            ))}
          </ul>
        )}
              <TextBox />

      </div>
     
      <style>{`
        .sidebar {
          padding: 16px;
          background-color: #f3f4f6;
          border-radius: 8px;
        }

        .search-container {
          display: flex;
          align-items: center;
          background-color: #1b4e5a;
          border-radius: 8px;
          padding: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }

        .search-input {
          flex: 1;
          border: none;
          outline: none;
          padding: 8px;
          font-size: 16px;
          border-radius: 8px;
          background-color: #1b4e5a; /* Explicitly set background color */
          color: white; /* Explicitly set text color */
        }

        .search-input::placeholder {
          color: rgba(255, 255, 255, 0.7); /* Set placeholder text color to light white */
          font-style: italic; /* Optional: set placeholder text style */
        }

        .search-input:focus {
          outline: none; /* Remove halo highlight */
          box-shadow: none; /* Remove any box shadow */
        }

        .search-icon {
          width: 24px;
          height: 24px;
          margin-left: 8px;
          margin-right: 100px;
          stroke: white; /* Set stroke color to white */
          z-index: 1; /* Ensure the icon is above other elements */
        }

        .search-results {
          margin-top: 8px;
          list-style: none;
          padding: 0;
          background-color: white;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          color: black;
        }

        .search-result-item {
          padding: 8px;
          border-bottom: 1px solid #dedcdc;
          cursor: pointer;
        }

        .search-result-item:last-child {
          border-bottom: none;
        }
      `}</style>
    </>
  );
};

export default SearchBar;