import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { apiService } from '../services/api';
import './LandingPage.css';

export const LandingPage: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (isLogin) {
        const loginValue = name || email;
        const tokenResponse = await apiService.authenticate({
          Login: loginValue,
          Password: password,
        });
        login(tokenResponse);
        navigate('/sheets');
      } else {
        await apiService.createUser({ Name: name, Email: email, Password: password });
        const tokenResponse = await apiService.authenticate({ Login: name, Password: password });
        login(tokenResponse);
        navigate('/sheets');
      }
    } catch (err) {
      const apiError = apiService.handleError(err);
      if (isLogin && apiError.status === 400) {
        setError('Неверный логин или пароль');
      } else {
        setError(apiError.message || 'Произошла ошибка');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="landing-page">
      <div className="landing-container">
        <div className="landing-header">
          <h1>SheetStruct</h1>
          <p className="version">v0.1.0</p>
        </div>

        <div className="auth-form-container">
          <div className="auth-tabs">
            <button
              className={isLogin ? 'active' : ''}
              onClick={() => setIsLogin(true)}
            >
              Вход
            </button>
            <button
              className={!isLogin ? 'active' : ''}
              onClick={() => setIsLogin(false)}
            >
              Регистрация
            </button>
          </div>

          <form onSubmit={handleSubmit} className="auth-form">
            {error && <div className="error-message">{error}</div>}

            {isLogin ? (
              <>
                <div className="form-group">
                  <label>Имя или Email</label>
                  <input
                    type="text"
                    value={name || email}
                    onChange={(e) => {
                      if (e.target.value.includes('@')) {
                        setEmail(e.target.value);
                        setName('');
                      } else {
                        setName(e.target.value);
                        setEmail('');
                      }
                    }}
                    placeholder="Введите имя или email"
                    required
                  />
                </div>
              </>
            ) : (
              <>
                <div className="form-group">
                  <label>Имя</label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Введите ваше имя"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Email</label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Введите ваш email"
                    required
                  />
                </div>
              </>
            )}

            <div className="form-group">
              <label>Пароль</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Введите пароль"
                required
              />
            </div>

            <button type="submit" className="submit-button" disabled={loading}>
              {loading ? 'Обработка...' : isLogin ? 'Войти' : 'Зарегистрироваться'}
            </button>
          </form>

          {isLogin && (
            <div className="reset-password-link">
              <button
                type="button"
                onClick={() => navigate('/reset-password')}
                className="link-button"
              >
                Забыли пароль?
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

