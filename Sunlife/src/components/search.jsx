import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import StockGraph from './stock-graph';

const Search = () => {
  const { stock } = useParams();
  const navigate = useNavigate();

  const handleHomeClick = () => {
    navigate('/');
  };

  return (
    <div className="stock-info">
      <div className="header">
        <button className="home-button" onClick={handleHomeClick}>
          Home
        </button>
        <h1>Stock Information</h1>
      </div>
      <div className="stock-details">
        <h2>{decodeURIComponent(stock)}</h2>
        <p>Here are some details about the stock {decodeURIComponent(stock)}.</p>
      </div>
      <div className="stock-graph">
        <StockGraph />
      </div>
      <style jsx>{`
        .stock-info {
          padding: 20px;
          background-color: var(--color-cornsilk);
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          max-width: 800px;
          margin: 40px auto;
          font-family: 'Gupter', serif;
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .header {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .stock-details {
          background-color: white;
          padding: 20px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }

        .home-button {
          background-color: var(--color-midnight-green);
          color: white;
          border: none;
          padding: 10px 20px;
          border-radius: 8px;
          cursor: pointer;
          font-size: 16px;
        }

        .home-button:hover {
          background-color: var(--color-jonquil);
        }

        h1 {
          color: var(--color-midnight-green);
          font-size: 2.5em;
          margin: 0;
        }

        h2 {
          color: var(--color-jonquil);
          font-size: 1.8em;
          margin-bottom: 10px;
        }

        p {
          color: #333;
          font-size: 1.2em;
        }

        .stock-graph {
          background-color: white;
          padding: 20px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
      `}</style>
    </div>
  );
};

export default Search;