import React, { useState, useRef, useEffect } from "react";
import { FiSend, FiMic, FiCopy } from "react-icons/fi";
import { format } from "date-fns";
import { useSymbol } from "../contexts/symbol-context";

const TextBox = () => {
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef(null);
  const {selectedSymbol, setSelectedSymbol} = useSymbol();
  const [isWaitingForResponse, setIsWaitingForResponse] = useState(false)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    if(selectedSymbol === '' || isWaitingForResponse) return;
    const newMessage = {
      text: `Can you tell me why ${selectedSymbol} has been performing like this recently?`,
      timestamp: new Date(),
      className: "userMessage"
    }
    setMessages([...messages, newMessage]);
    setSelectedSymbol('');
  }, [selectedSymbol])

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    if(messages.length === 0 || (messages.slice(-1)).className === 'aiMessage') return;
    setIsWaitingForResponse(true);
  }, [messages])

  const handleSendMessage = () => {
    if (inputMessage.trim() !== "") {
      const newMessage = {
        text: inputMessage,
        timestamp: new Date(),
        className: "userMessage"
      };
      setMessages([...messages, newMessage]);
      setInputMessage("");
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <div className="chat-container">
      <div className="messages">
        {messages.map((message, index) => (
          <div key={index} className={message.className}>
            <span>{message.text}</span>
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
            {/* <textarea
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Type your message..."
              className="input-message"
            /> */}
            {isWaitingForResponse ? (
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
              onKeyPress={handleKeyPress}
              placeholder="Type your message..."
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

        .userMessage {
          background-color: var(--color-cornsilk);
          padding: 8px;
          border-radius: 8px;
          margin-bottom: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          color: black;
          margin-left: 10%;
        }

        .aiMessage {
          background-color: white;
          padding: 8px;
          border-radius: 8px;
          margin-bottom: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          color: black;
          margin-right: 10%;
        }

        .timestamp {
          display: block;
          font-size: 0.8em;
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
      `}</style>
    </div>
  );
};

export default TextBox;