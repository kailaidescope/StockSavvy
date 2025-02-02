import React, { useState, useEffect } from 'react';
import { getTickerInfo } from '../data/api-requests.js';
import "./EmojiScrollbar.css";

export default function EmojiScrollbar({ symbol, emojiTop, emojiBottom, titletop, titlebottom }) {
  const [sentimentValue, setSentimentValue] = useState(50);
  const [numArticles, setNumArticles] = useState(0);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      const cachedData = localStorage.getItem(`tickerInfo-${symbol}`);
      if (cachedData) {
        const parsed = JSON.parse(cachedData);
        setSentimentValue(parsed.sentimentValue);
        setNumArticles(parsed.numArticles);

        setLoading(false);
        return;
      }
      try {
        const infoString = await getTickerInfo(symbol);
        const newsData = JSON.parse(infoString);
        
        
        setSentimentValue(newsData.avg_sentiment );
        setNumArticles(newsData.num_articles/100);
        
        localStorage.setItem(`tickerInfo-${symbol}`, JSON.stringify({
          sentimentValue: normalizedSentiment,
          numArticles: newsData.num_articles
        }));
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [symbol]);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="emoji-scroll-container">
      <span className="emoji-top" title={titletop}>{emojiTop}</span>
      <input
        type="range"
        min="0"
        max="100"
        value={titletop === "Hot topic" ? numArticles : sentimentValue}
        className="emoji-range"
        onChange={(e) => setValue(e.target.value)}
      />
      <span className="emoji-bottom"title={titlebottom}>{emojiBottom}</span>
    </div>
  );
}