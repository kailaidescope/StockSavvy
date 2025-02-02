import axios from 'axios';
import fake_data from '../data/fake_investment_data.js';

const BASE_URL = 'http://localhost:3000/api/v1';

// Configuration for axios requests
const axiosConfig = {
    headers: {
        'Content-Type': 'application/json'
    }
};

/**
 * Get historical data for a specific ticker
 * @param {string} symbol - The stock symbol
 * @returns {Promise} - Returns historical price data
 */
export const getTickerHistory = async (symbol) => {
    try {
        const response = await axios.get(`${BASE_URL}/stocks/tickers/${symbol}/history`, axiosConfig);
        return JSON.stringify(response.data.history);
    } catch (error) {
        console.error('Error fetching ticker history:', error);
        return fake_data; // Fallback to fake data if API call fails
    }
};

/**
 * Get news sentiment for a specific ticker
 * @param {string} symbol - The stock symbol
 * @returns {Promise} - Returns news sentiment data
 */
export const getTickerNews = async (symbol) => {
    try {
        const response = await axios.get(`${BASE_URL}/stocks/tickers/${symbol}/news`, axiosConfig);
        return JSON.stringify(response.data);
    } catch (error) {
        console.error('Error fetching ticker news:', error);
        return []; // Return empty array if API call fails
    }
};

/**
 * Get basic information about a ticker
 * @param {string} symbol - The stock symbol
 * @returns {Promise} - Returns ticker information
 */
export const getTickerInfo = async (symbol) => {
    try {
        const response = await axios.get(`${BASE_URL}/stocks/tickers/${symbol}`, axiosConfig);
        return JSON.stringify(response.data);
    } catch (error) {
        console.error('Error fetching ticker info:', error);
        return [];
    }
};

/**
 * Get user's stock holdings
 * @returns {Promise} - Returns holdings information
 */
export const getHoldings = async () => {
    try {
        const response = await axios.get(`${BASE_URL}/stocks/holdings`, axiosConfig);
        return JSON.stringify(response.data);
    } catch (error) {
        console.error('Error fetching holdings:', error);
        return [];
    }
};
