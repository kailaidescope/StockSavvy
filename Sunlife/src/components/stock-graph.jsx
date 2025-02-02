import React, { useEffect, useRef, useState } from 'react';
import { createChart, AreaSeries } from 'lightweight-charts';
import { getTickerHistory } from '../data/api-requests';

const StockGraph = ({ symbol = 'AAPL' }) => {
    const chartContainerRef = useRef(null);
    const [data, setData] = useState([]);
    
    const fetchData = async () => {
        try {
            const response = await getTickerHistory(symbol);
            const parsedData = JSON.parse(response);
            setData(parsedData);
            console.log('Fetched data:', parsedData);
        } catch (error) {
            console.error('Error:', error);
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

            const lineSeries = chart.addSeries(AreaSeries, {
                color: '#1b4e5a',
                lineWidth: 3
            });

            lineSeries.setData(data);
            chart.timeScale().fitContent();

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
        <div>
            <button className="refresh-button" onClick={fetchData}>Refresh</button>
            <div 
                ref={chartContainerRef} 
                style={{ 
                    width: '100%', 
                    borderRadius: '15px', // Add rounded corners
                    overflow: 'hidden', // Ensure content doesn't overflow the rounded corners
                    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)', // Add box shadow
                    backgroundColor: 'white', // Ensure background color is white
                    padding: '16px' // Add some padding
                }} 
            />
            <style jsx>{`
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