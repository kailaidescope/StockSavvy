import React from 'react';

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
        </div>
      </div>
      <style jsx>{`
        .sidebar {
          padding: 16px;
          background-color: var(--color-platinum);
          border-radius: 8px;
          box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }

        .search-container {
          display: flex;
          align-items: center;
          background-color: white;
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
        }

        .search-icon {
          width: 24px;
          height: 24px;
          margin-left: 8px;
          color: var(--color-midnight-green);
        }
      `}</style>
    </>
  );
};

export default SearchBar;