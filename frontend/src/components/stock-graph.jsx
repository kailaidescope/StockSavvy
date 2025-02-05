import React, { useEffect, useRef, useState } from 'react';
import { createChart, AreaSeries, LineSeries } from 'lightweight-charts';
import { getHoldingHistory, getTickerHistory } from '../data/api-requests';

const StockGraph = ({ symbol = 'AAPL' , isHolding = false}) => {
    const chartContainerRef = useRef(null);
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchData = async () => {
        setLoading(true);
        //localStorage.clear();
        const cachedData = localStorage.getItem(`stockData-${symbol}`);
        if (isHolding) {
            try {
                // const response = await getTickerHistory(symbol);
                const response = await getHoldingHistory(symbol);
                // const response = await getHoldings();
                setData(response);
                localStorage.setItem(`stockData-${symbol}`, JSON.stringify(response));
                setLoading(false);
            } catch (error) {
                console.error('Error:', error);
                setLoading(false);
            }
        }
        else if (cachedData) {
            setData(JSON.parse(cachedData));
            setLoading(false);
        } 
    };
    useEffect(() => {
        fetchData();
    }, [symbol]);

    useEffect(() => {
        if (chartContainerRef.current && data.length > 0) {
            const chart = createChart(chartContainerRef.current, {
                width: chartContainerRef.current.clientWidth,
                height: 250,
                layout: {
                    textColor: 'black',
                    background: { type: 'solid', color: 'white' }
                }
            });

            const valueSeries = chart.addSeries(AreaSeries, {
                color: '#1b4e5a',
                lineWidth: 3
            });

            valueSeries.setData(data);
            chart.timeScale().fitContent();

            const stockAmtSeries = chart.addSeries(LineSeries, {
                color: '#f5a623',
                lineWidth: 2
            });

            const stockAmtData = data.map((d) => ({
                time: d.time,
                value: d.shares * 100
            }));
            stockAmtSeries.setData(stockAmtData);
            

            // Handle resize
            const handleResize = () => {
                chart.applyOptions({
                    width: chartContainerRef.current.clientWidth
                });
            };

            window.addEventListener('resize', handleResize);

            return () => {
                window.removeEventListener('resize', handleResize);
                chart.remove();
            };
        }
    }, [data]);

    return (
        <div className="chart-wrapper">
            {loading ? (
                <div className="loading-container">
                    <div className="loading-bar">
                        <div className="loading-progress"></div>
                    </div>
                    <p>Loading {symbol} data...</p>
                </div>
            ) : (
                <div ref={chartContainerRef} 
                style={{ 
                    width: '100%', 
                    borderRadius: '15px', // Add rounded corners
                    overflow: 'hidden', // Ensure content doesn't overflow the rounded corners
                    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)', // Add box shadow
                    backgroundColor: 'white', // Ensure background color is white
                    padding: '16px' // Add some padding
                }} 
            />)}
            <style jsx>{`
             .chart-wrapper {
                    position: relative;
                    height: 250px;
                    width: 100%;
                    background: white;
                    border-radius: 8px;
                }
                
                .loading-container {
                    position: absolute;
                    top: 50%;
                    left: 50%;
                    transform: translate(-50%, -50%);
                    text-align: center;
                }
                
                .loading-bar {
                    width: 200px;
                    height: 4px;
                    background: #f0f0f0;
                    border-radius: 2px;
                    overflow: hidden;
                    margin-bottom: 12px;
                }
                
                .loading-progress {
                    width: 40%;
                    height: 100%;
                    background: #1b4e5a;
                    border-radius: 2px;
                    animation: loading 1.5s infinite ease-in-out;
                }
                
                @keyframes loading {
                    0% {
                        transform: translateX(-250%);
                    }
                    100% {
                        transform: translateX(250%);
                    }
                }
                
                p {
                    color: #666;
                    font-size: 14px;
                    margin: 0;
                    font-weight: 500;
                }
                .refresh-button {
                    background-color: var(--color-midnight-green);
                    color: white;
                    border: none;
                    padding: 10px 20px;
                    border-radius: 8px;
                    cursor: pointer;
                    font-size: 16px;
                    margin-bottom: 20px;
                }

                .refresh-button:hover {
                    background-color: var(--color-jonquil);
                }
            `}</style>
        </div>
    );
};

export default StockGraph;