import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import InvestmentDashboard from './components/investment-dashboard';
import Search from './components/search-page';
import { SymbolProvider } from './contexts/symbol-context';

const App = () => {
  return (
    <SymbolProvider>
      <Router>
        <Routes>
          <Route path="/" element={<InvestmentDashboard />} />
          <Route path="/search/:stock" element={<Search />} />
        </Routes>
      </Router>
    </SymbolProvider> 
  );
};

export default App;