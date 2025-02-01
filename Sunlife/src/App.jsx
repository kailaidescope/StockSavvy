import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import InvestmentDashboard from './components/investment-dashboard';
import Search from './components/search';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<InvestmentDashboard />} />
        <Route path="/search/:stock" element={<Search />} />
      </Routes>
    </Router>
  );
};

export default App;