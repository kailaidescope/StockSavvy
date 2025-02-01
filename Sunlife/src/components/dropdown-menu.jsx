import React, { useState } from 'react';

const DropdownMenu = () => {
  const [isOpen, setIsOpen] = useState(false);

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  return (
    <div className="dropdown">
      <button onClick={toggleDropdown} className="dropdown-toggle">
        Menu
      </button>
      {isOpen && (
        <div className="dropdown-menu">
          <a href="#" className="dropdown-item">Log In</a>
          <a href="#" className="dropdown-item">Sign Up</a>
          <a href="#" className="dropdown-item">Profile</a>
        </div>
      )}
      <style jsx>{`
        .dropdown {
          position: relative;
          display: inline-block;
        background-color: #ffcb05;

        }

        .dropdown-toggle {
          background-color: #ffcb05;
          color: white;
          padding: 8px 16px;
          border: none;
          cursor: pointer;
        }

        .dropdown-menu {
          display: block;
          position: absolute;
          right: 0;
          background-color: white;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          overflow: hidden;
          z-index: 1;
        }

        .dropdown-item {
          padding: 8px 16px;
          text-decoration: none;
          color: #1b4e5a;
          display: block;
        }

        .dropdown-item:hover {
          background-color: #f3f4f6;
        }
      `}</style>
    </div>
  );
};

export default DropdownMenu;