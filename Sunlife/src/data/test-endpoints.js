import { 
    getTickerHistory, 
    getTickerNews, 
    getTickerInfo,
    getHoldings
} from './api-requests.js';

const testSymbol = 'AAPL';

const testEndpoints = async () => {
    console.log('Starting API endpoint tests...\n');

    try {
        // Test getTickerHistory
        console.log('*************Testing getTickerHistory...**********');
        const history = await getTickerHistory(testSymbol);
        console.log('History data:', history, '\n');

        // Test getTickerNews
        console.log('************Testing getTickerNews...**********');
        const news = await getTickerNews(testSymbol);
        console.log('News data:', news, '\n');

        // Test getTickerInfo
        console.log('************Testing getTickerInfo...***********');
        const info = await getTickerInfo(testSymbol);
        console.log('Info data:', info, '\n');

        // Test getHoldings
        console.log('**********Testing getHoldings...**********');
        const holdings = await getHoldings();
        console.log('Holdings data:', holdings, '\n');

    } catch (error) {
        console.error('💀Test failed:', error);
    }
};

// Run tests
testEndpoints();

// To run this file:
// node src/tests/api-tests.js