import React, { useState, useRef, useEffect } from "react";
import { FiSend, FiMic, FiCopy } from "react-icons/fi";
import { format } from "date-fns";
import { useSymbol } from "../contexts/symbol-context";
import { useLocation } from "react-router-dom";
import { sendChat } from "../data/api-requests";
import ReactMarkdown from "react-markdown";

const TextBox = () => {
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef(null);
  const { selectedSymbol, setSelectedSymbol, selectedSymbols, setSelectedSymbols } = useSymbol();
  const [isWaitingForResponse, setIsWaitingForResponse] = useState(false)
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };
  const location = useLocation()

  const stockData = [
    { symbol: 'AAPL', name: 'Apple Inc.', sector: 'Technology', price: '180.95', change: '+1.2%' },
    { symbol: 'MSFT', name: 'Microsoft Corp.', sector: 'Technology', price: '378.85', change: '+0.8%' },
    { symbol: 'JNJ', name: 'Johnson & Johnson', sector: 'Healthcare', price: '155.42', change: '-0.5%' },
    { symbol: 'PFE', name: 'Pfizer Inc.', sector: 'Healthcare', price: '28.79', change: '+1.1%' },
    { symbol: 'JPM', name: 'JPMorgan Chase', sector: 'Finance', price: '167.42', change: '+0.3%' },
    { symbol: 'BAC', name: 'Bank of America', sector: 'Finance', price: '33.98', change: '-0.7%' }
  ];

  const sectors = [...new Set(stockData.map(stock => stock.sector))];

  useEffect(() => {
    if (selectedSymbols.length === 0) return;
    let text = `Can you tell me more about $${selectedSymbols.toString()}?`
    text = text.replaceAll(',', ', $');
    setInputMessage(text)
  }, [selectedSymbols])

  useEffect(() => {
    if (selectedSymbol === '' || isWaitingForResponse) return;
    const newMessage = {
      text: `Can you tell me why $${selectedSymbol} has been performing like this recently?`,
      timestamp: Math.floor(Date.now() / 1000),
      sender: "user"
    }
    setMessages((prevMessages) => [...prevMessages, newMessage]);
    setSelectedSymbol('');
  }, [selectedSymbol])

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    if (messages.length === 0 || (messages.slice(-1)).sender === 'ai') return;
    setIsWaitingForResponse(true);
  }, [messages])

  const handleSendMessage = () => {
    if (inputMessage.trim() !== "") {
      const newMessage = {
        text: inputMessage,
        timestamp: Math.floor(Date.now() / 1000),
        sender: "user"
      };
      setMessages((prevMessages) => [...prevMessages, newMessage]);
      setInputMessage("");
      setSelectedSymbols([]);
      setIsTyping(true); // Show typing indicator
      sendChat(newMessage.text, messages).then(response => {
        const newMessage = {
          text: response,
          timestamp: Math.floor(Date.now() / 1000),
          sender: "ai"
        };
        setMessages((prevMessages) => [...prevMessages, newMessage]);
        setIsTyping(false); // Hide typing indicator
        setIsWaitingForResponse(false);
      });
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
    else if (e.key === "Delete" || e.key === "Backspace") {
      setSelectedSymbols([]);
    }
  };

  const handleSectorButtonPress = (sector) => {
    setSelectedSymbols([]);
    setInputMessage(`Can you tell me more about the current financial state of the ${sector} sector?`);
  }

  return (
    <div className="chat-container">
      {location.pathname === '/advanced-search' ? (
        <div className="sector-buttons-container">
          <h2 className="sector-buttons-title">Learn more about...</h2>
          <div className="sector-buttons">
            {sectors.map(sector => (
              <button
                key={sector}
                className={`sector-button`}
                onClick={() => handleSectorButtonPress(sector.charAt(0).toUpperCase() + sector.slice(1))}
              >
                {sector.charAt(0).toUpperCase() + sector.slice(1)}
              </button>
            ))}
          </div>
        </div>
      ) : (null)}
      <div className="messages">
        {messages.map((message, index) => (
          <div key={index} className={message.sender}>
            <ReactMarkdown>{message.text}</ReactMarkdown> {/* Markdown Rendering */}
            <span className="timestamp">
              {format(message.timestamp, "p, MMMM dd")}
            </span>
          </div>
        ))}
        {isTyping && (
          <div className="typing-indicator">
            <div className="dot"></div>
            <div className="dot delay-100"></div>
            <div className="dot delay-200"></div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input Section */}
      <div className="input-section">
        <div className="input-container">
          <button className="icon-button" aria-label="Voice input">
            <FiMic size={20} />
          </button>
          <div className="textarea-container">
            {false ? (
              <textarea
                disabled
                value={inputMessage}
                placeholder="Waiting for response..."
                className="input-message"
              />
            ) : (
              <textarea
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyDown={handleKeyPress}
                placeholder="Type message here..."
                className="input-message"
              />
            )}
          </div>
          <button className="icon-button" onClick={handleSendMessage}>
            <FiSend size={20} />
          </button>
          <button className="icon-button">
            <FiCopy size={20} />
          </button>
        </div>
      </div>

      <style jsx>{`
        React-Markdown, p {
          margin-bottom: 8px;
          color: black;
          margin-right: 10%;
          font-size: 1.1em;
        }
      
        .chat-container {
          display: flex;
          flex-direction: column;
          height: 90%;
          width: 100%;
          background-color: #f3f4f6;
          border-radius: 8px;
          padding: 16px;
          box-sizing: border-box;
        }

        .messages {
          flex: 1;
          overflow-y: auto;
          margin-bottom: 16px;
        }

        .user {
          background-color: var(--color-cornsilk);
          padding: 8px;
          border-radius: 8px;
          margin-bottom: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          color: black;
          margin-left: 10%;
          font-size: 1.5em;
        }

        .ai {
          background-color: white;
          padding: 8px;
          border-radius: 8px;
          margin-bottom: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          color: black;
          margin-right: 10%;
          font-size: 1.5em;
        }

        .timestamp {
          display: block;
          font-size: 0.6em;
          color: #888;
          margin-top: 4px;
        }

        .typing-indicator {
          display: flex;
          align-items: center;
          gap: 4px;
          color: #888;
        }

        .dot {
          width: 8px;
          height: 8px;
          background-color: #888;
          border-radius: 50%;
          animation: bounce 1s infinite;
        }

        .dot.delay-100 {
          animation-delay: 0.1s;
        }

        .dot.delay-200 {
          animation-delay: 0.2s;
        }

        @keyframes bounce {
          0%, 100% {
            transform: translateY(0);
          }
          50% {
            transform: translateY(-8px);
          }
        }

        .input-section {
          padding: 16px;
          border-top: 1px solid #ddd;
          background-color: white;
          border-radius: 10px
        }

        .input-container {
          display: flex;
          align-items: flex-end;
          gap: 8px;
        }

        .icon-button {
          background-color: transparent;
          border: none;
          color: #888;
          cursor: pointer;
          transition: color 0.3s;
        }

        .icon-button:hover {
          color: #555;
        }

        .textarea-container {
          flex: 1;
          position: relative;
        }

        .input-message {
          width: 100%;
          padding: 8px;
          font-size: 16px;
          border-radius: 8px;
          border: 1px solid #ccc;
          resize: none;
          box-sizing: border-box;
          background-color: white;
          color: black;
        }
        .sector-buttons-container {
          display: flex;
          flex-direction: column;
          gap: 10px;
          margin-top: 15px;
          margin-bottom: 15px;
        }
        .sector-button {
          padding: 8px 16px;
          border: none;
          border-radius: 20px;
          background-color: var(--color-platinum);
          color: var(--color-midnight-green);
          cursor: pointer;
          transition: all 0.3s ease;
        }

        .sector-buttons {
          display: flex;
          gap: 10px;
          flex-wrap: wrap;
        }
        .sector-buttons-title {
          color: black;
        }
      `}</style>
    </div>
  );
};

export default TextBox;