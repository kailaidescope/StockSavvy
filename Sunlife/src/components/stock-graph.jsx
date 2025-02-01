import React, { useEffect, useRef, useState } from 'react';
import { createChart, LineSeries } from 'lightweight-charts';
import axios from 'axios';
import fake_data from '../data/fake_investment_data'

const StockGraph = () => {
    const chartContainerRef = useRef(null);
    const [data, setData] = useState([]);
    
    let config = {
        method: 'get',
        maxBodyLength: Infinity,
        url: 'http://172.20.10.2:3000/api/v1/stocks/tickers/hi/history',
        headers: { }
      };
      
    const fetchData = () => {
        axios.request(config)
            .then(response => {
                setData(response.data);
                console.log(JSON.stringify(response.data));
            })
            .catch(error => {
                setData(fake_data)
                console.error('Error fetching data:', error);
            });
    };

    useEffect(() => {
        fetchData();
    }, []);

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

            const lineSeries = chart.addSeries(LineSeries, {
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