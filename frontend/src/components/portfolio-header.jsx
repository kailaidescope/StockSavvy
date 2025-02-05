import React, { useState, useEffect, useMemo } from "react";
import { FaArrowUp, FaArrowDown } from "react-icons/fa";
import CountUp from "react-countup";

const PortfolioHeader = () => {
  const [investmentData, setInvestmentData] = useState({
    totalHoldings: 156789.42,
    percentageChange: 2.35,
    changeDirection: "positive",
    isLoading: false,
    error: null
  });

  const formatCurrency = (value) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD"
    }).format(value);
  };

  const colorClass = investmentData.changeDirection === "positive" ? "positive" : "negative";
  const Arrow = investmentData.changeDirection === "positive" ? FaArrowUp : FaArrowDown;

  return (
    <div className="portfolio-header">
      <div className="holdings">
        <span>Holdings: {formatCurrency(investmentData.totalHoldings)}</span>
      </div>
      <div className={`percentage-change ${colorClass}`} role="status" aria-label={`${investmentData.percentageChange}% ${investmentData.changeDirection === "positive" ? "increase" : "decrease"}`}>
        <Arrow className="arrow-icon" />
        <CountUp
          start={0}
          end={investmentData.percentageChange}
          duration={2}
          decimals={2}
          suffix="%"
        />
      </div>
      <style jsx>{`
        .portfolio-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 20px;
          background-color: var(--color-cornsilk);
          border-radius: 8px;
          box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
          margin-bottom: 20px;
          transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .portfolio-header:hover {
          transform: translateY(-5px);
          box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
        }

        .holdings {
          font-size: 1.2em;
          color: var(--color-midnight-green);
        }

        .holdings span {
          font-weight: bold;
          font-size: 1.5em;
        }

        .percentage-change {
          display: flex;
          align-items: center;
          font-size: 1.2em;
          font-weight: 600;
        }

        .percentage-change.positive {
          color: green;
        }

        .percentage-change.negative {
          color: red;
        }

        .arrow-icon {
          font-size: 1em;
          margin-right: 0.5em;
        }
      `}</style>
    </div>
  );
};

export default PortfolioHeader;