import React, { useEffect, useRef } from 'react';
import { createChart, LineSeries } from 'lightweight-charts';
import data from '../data/fake_investment_data.js';

const StockGraph = () => {
    const chartContainerRef = useRef(null);

    useEffect(() => {
        if (chartContainerRef.current) {
            const chart = createChart(chartContainerRef.current, {
                width: chartContainerRef.current.clientWidth,
                height: 250,
                layout: {
                    textColor: 'black',
                    background: { type: 'solid', color: 'white' }
                }
            });

            const lineSeries = chart.addSeries(LineSeries, {
                color: '#2962FF',
                lineWidth: 2
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
    }, []);

    return (
        <div 
            ref={chartContainerRef} 
            style={{ 
                width: '100%', 
                height: '100%',
                minHeight: '400px' 
            }} 
        />
    );
};

export default StockGraph;