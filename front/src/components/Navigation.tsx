import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Navigation.css';

export const Navigation: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const location = useLocation();

  if (!isAuthenticated) {
    return null;
  }

  return (
    <nav className="navigation">
      <div className="nav-container">
        <Link to="/sheets" className="nav-logo">
          SheetStruct
        </Link>
        <div className="nav-links">
          <Link
            to="/sheets"
            className={location.pathname.startsWith('/sheet') ? 'active' : ''}
          >
            Листы
          </Link>
          <Link
            to="/settings"
            className={location.pathname === '/settings' ? 'active' : ''}
          >
            Настройки
          </Link>
        </div>
        <div className="nav-user">
          <span className="user-name">{user?.Name}</span>
          <button onClick={logout} className="logout-button">
            Выход
          </button>
        </div>
      </div>
    </nav>
  );
};

