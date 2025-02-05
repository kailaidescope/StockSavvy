import React, { useEffect, useRef, useState } from 'react';
import { createChart, AreaSeries } from 'lightweight-charts';
import data from '../data/fake_investment_data'

const TickerGraph = () => {
    const chartContainerRef = useRef(null);

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
        </div>
    );
};

export default TickerGraph;