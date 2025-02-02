import React, { useState, useEffect } from 'react';
import { getTickerInfo, getTickerNews } from '../data/api-requests.js';
import "./EmojiScrollbar.css";

export default function EmojiScrollbar({ symbol, emojiTop, emojiBottom, titletop, titlebottom }) {
  const [sentimentValue, setSentimentValue] = useState(50);
  const [numArticles, setNumArticles] = useState(0);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      const cachedData = localStorage.getItem(`tickerInfo-${symbol}`);
      //const cachedData = null;
      if (cachedData) {
        const parsed = JSON.parse(cachedData);
        setSentimentValue(parsed.sentimentValue);
        setNumArticles(parsed.numArticles);

        setLoading(false);
        return;
      }
      try {
        const infoString = await getTickerNews(symbol);
        const newsData = JSON.parse(infoString);

        //console.log(newsData);
        
        
        setSentimentValue(newsData.avg_sentiment);
        setNumArticles(newsData.num_articles);
        
        localStorage.setItem(`tickerInfo-${symbol}`, JSON.stringify({
          sentimentValue: newsData.avg_sentiment,
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
        step="0.01"
        min={titletop === "Hot topic" ? 0 : 0}
        max={titletop === "Hot topic" ? 350 : 2}
        value={titletop === "Hot topic" ? numArticles : sentimentValue+1}
        className="emoji-range"
        onChange={(e) => {
          const value = e.target.value;
          if (titletop === "Hot topic") {
            setNumArticles(value);
          } else {
            setSentimentValue(value);
          }
        }}
      />
      <span className="emoji-bottom" title={titlebottom}>{emojiBottom}</span>
    </div>
  );
}