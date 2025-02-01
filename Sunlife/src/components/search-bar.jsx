import React from 'react';
import DropdownMenu from './dropdown-menu';

const SearchBar = () => {
  return (
    <>
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
          <DropdownMenu />

        </div>
      </div>
      <style jsx>{`
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

        .search-icon {
          width: 24px;
          height: 24px;
          margin-left: 8px;
          stroke: white;
        }
      `}</style>
    </>
  );
};

export default SearchBar;