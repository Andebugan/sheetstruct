import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiService } from '../services/api';
import './ResetPasswordPage.css';

export const ResetPasswordPage: React.FC = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await apiService.resetPassword({
        ...(name ? { Name: name } : { Email: email }),
      });
      setSuccess(true);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="reset-password-page">
        <div className="reset-container">
          <h1>Сброс пароля</h1>
          <p>Если аккаунт с указанными данными существует, письмо для сброса пароля было отправлено.</p>
          <button onClick={() => navigate('/login')} className="back-button">
            Вернуться к входу
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="reset-password-page">
      <div className="reset-container">
        <h1>Сброс пароля</h1>
        <form onSubmit={handleSubmit} className="reset-form">
          {error && <div className="error-message">{error}</div>}

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

          <button type="submit" className="submit-button" disabled={loading}>
            {loading ? 'Отправка...' : 'Отправить письмо для сброса'}
          </button>

          <button
            type="button"
            onClick={() => navigate('/login')}
            className="back-link"
          >
            Вернуться к входу
          </button>
        </form>
      </div>
    </div>
  );
};

