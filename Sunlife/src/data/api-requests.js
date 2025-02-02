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
        return response.data.history;
    } catch (error) {
        console.error('Error fetching ticker history:', error);
        return []; // Fallback to fake data if API call fails
    }
};

/**
 * Get historical data for a specific holding
 * @param {string} symbol - The stock symbol
 * @returns {Promise} - Returns historical value of your holding
 */
export const getHoldingHistory = async (symbol) => {
    try {
        const response = await axios.get(`${BASE_URL}/stocks/holdings/${symbol}`, axiosConfig);
        return response.data.history;
    } catch (error) {
        console.error('Error fetching holding history:', error);
        return []; // Fallback to fake data if API call fails
    }
}

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

/**
 * Send a message to the chatbot and receive a response
 * @param {string} message - The user's message
 * @returns {Promise} - Returns the chatbot's response
 */
export const sendChat = async (newMessageText, messages) => {
    let data;
    // Try to encode chat message
    try {
        console.log(newMessageText, "\n",messages)
        data = JSON.stringify({
            "prompt": newMessageText,
            "history": messages
          });
        console.log(data)
    } catch (error) {
        console.error('Error encoding chat message:', error);
        return [];
    }
    
    // Try to send request
    try {
        let chatRequestConfig = {
            method: 'post',
            maxBodyLength: Infinity,
            url: `${BASE_URL}/chat`,
            headers: { 
              'Content-Type': 'application/json'
            },
            data : data
          };
        const response = await axios.request(chatRequestConfig)
        let output = response.data;
        let textOutput = output["ai-response"];
        console.log(output, "\n", textOutput);
        return textOutput;
    } catch (error) {
        console.error('Error sending chat:', error);
        return "I'm sorry, I don't understand that.";
    }
}
