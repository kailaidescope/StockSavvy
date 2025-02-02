import React, { useState } from 'react';
import TextBox from './text-box';
import { useNavigate } from 'react-router-dom';
import { useSymbol } from '../contexts/symbol-context';

const AdvancedSearchPage = () => {
    const navigate = useNavigate();
    const [selectedSector, setSelectedSector] = useState('all');
    const [searchTerm, setSearchTerm] = useState('');
    const { selectedSymbols, setSelectedSymbols } = useSymbol();

    const stockData = [
        { symbol: 'AAPL', name: 'Apple Inc.', sector: 'Technology', price: '180.95', change: '+1.2%' },
        { symbol: 'MSFT', name: 'Microsoft Corp.', sector: 'Technology', price: '378.85', change: '+0.8%' },
        { symbol: 'JNJ', name: 'Johnson & Johnson', sector: 'Healthcare', price: '155.42', change: '-0.5%' },
        { symbol: 'PFE', name: 'Pfizer Inc.', sector: 'Healthcare', price: '28.79', change: '+1.1%' },
        { symbol: 'JPM', name: 'JPMorgan Chase', sector: 'Finance', price: '167.42', change: '+0.3%' },
        { symbol: 'BAC', name: 'Bank of America', sector: 'Finance', price: '33.98', change: '-0.7%' }
    ];

    const sectors = ['all', ...new Set(stockData.map(stock => stock.sector))];

    const filteredStocks = stockData.filter(stock => 
        (selectedSector === 'all' || stock.sector === selectedSector) &&
        (stock.symbol.toLowerCase().includes(searchTerm.toLowerCase()) || 
         stock.name.toLowerCase().includes(searchTerm.toLowerCase()))
    );

    const handleDragStart = (event, stock) => {
        event.dataTransfer.setData('text/plain', stock.symbol);
    };

    const handleDrop = (event) => {
        event.preventDefault();
        const stockSymbol = event.dataTransfer.getData('text/plain');
        if(selectedSymbols.includes(stockSymbol)) return;
        setSelectedSymbols([...selectedSymbols, stockSymbol]);
    };

    const handleDragOver = (event) => {
        event.preventDefault();
    };

    return (
        <div className="advanced-search-container">
            <div className="stock-section">
                <div className="search-filters">
                    <div className="search-header">
                        <button className="home-button" onClick={() => navigate('/')}>
                            Home
                        </button>
                        <input
                            type="text"
                            placeholder="Search stocks..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="search-input"
                        />
                    </div>
                    <div className="sector-filters">
                        {sectors.map(sector => (
                            <button
                                key={sector}
                                className={`sector-button ${selectedSector === sector ? 'active' : ''}`}
                                onClick={() => setSelectedSector(sector)}
                            >
                                {sector.charAt(0).toUpperCase() + sector.slice(1)}
                            </button>
                        ))}
                    </div>
                </div>
                <div className="stocks-list">
                    {filteredStocks.map(stock => (
                        <div 
                            key={stock.symbol} 
                            className="stock-card" 
                            draggable 
                            onDragStart={(event) => handleDragStart(event, stock)}
                        >
                            <div className="stock-info">
                                <h3>{stock.symbol}</h3>
                                <p>{stock.name}</p>
                            </div>
                            <div className="stock-details">
                                <span className="stock-price">${stock.price}</span>
                                <span className={`stock-change ${stock.change.startsWith('+') ? 'positive' : 'negative'}`}>
                                    {stock.change}
                                </span>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
            <div 
                className="chat-section" 
                onDrop={handleDrop} 
                onDragOver={handleDragOver}
            >
                <TextBox />
            </div>

            <style jsx>{`
                .advanced-search-container {
                    display: flex;
                    height: 100vh;
                    background-color: var(--color-cornsilk);
                    padding: 20px;
                    gap: 20px;
                }

                .home-button {
                    padding: 12px 20px;
                    background-color: var(--color-midnight-green);
                    color: white;
                    border: none;
                    border-radius: 8px;
                    cursor: pointer;
                    font-size: 16px;
                    transition: all 0.3s ease;
                    height: 100%;
                }

                .home-button:hover {
                    background-color: var(--color-jonquil);
                    transform: translateY(-2px);
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
                }

                .stock-section {
                    flex: 2;
                    display: flex;
                    flex-direction: column;
                    gap: 20px;
                }

                .chat-section {
                    flex: 1;
                    background: white;
                    border-radius: 12px;
                    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
                }

                .search-filters {
                    background: white;
                    padding: 20px;
                    border-radius: 12px;
                    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
                }

                .search-header {
                    display: flex;
                    gap: 15px;
                    align-items: center;
                    margin-bottom: 15px;
                }

                .search-input {
                    width: 100%;
                    padding: 12px;
                    border: 2px solid var(--color-platinum);
                    border-radius: 8px;
                    font-size: 16px;
                    transition: all 0.3s ease;
                }

                .search-input:focus {
                    outline: none;
                    border-color: var(--color-midnight-green);
                    box-shadow: 0 0 0 3px rgba(27, 78, 90, 0.2);
                }

                .sector-filters {
                    display: flex;
                    gap: 10px;
                    margin-top: 15px;
                    flex-wrap: wrap;
                }

                .sector-button {
                    padding: 8px 16px;
                    border: none;
                    border-radius: 20px;
                    background-color: var(--color-platinum);
                    color: var(--color-midnight-green);
                    cursor: pointer;
                    transition: all 0.3s ease;
                }

                .sector-button:hover {
                    background-color: var(--color-jonquil);
                    transform: translateY(-2px);
                }

                .sector-button.active {
                    background-color: var(--color-midnight-green);
                    color: white;
                }

                .stocks-list {
                    display: flex;
                    flex-direction: column;
                    gap: 10px;
                    overflow-y: auto;
                    padding: 10px;
                }

                .stock-card {
                    background: white;
                    padding: 15px;
                    border-radius: 12px;
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
                    transition: all 0.3s ease;
                }

                .stock-card:hover {
                    transform: translateY(-2px);
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
                }

                .stock-info h3 {
                    color: var(--color-midnight-green);
                    margin: 0;
                }

                .stock-info p {
                    color: #666;
                    margin: 5px 0 0 0;
                }

                .stock-details {
                    text-align: right;
                }

                .stock-price {
                    display: block;
                    font-size: 1.2em;
                    font-weight: bold;
                    color: var(--color-midnight-green);
                }

                .stock-change {
                    font-size: 0.9em;
                }

                .stock-change.positive {
                    color: #22c55e;
                }

                .stock-change.negative {
                    color: #ef4444;
                }
            `}</style>
        </div>
    );
};

export default AdvancedSearchPage;