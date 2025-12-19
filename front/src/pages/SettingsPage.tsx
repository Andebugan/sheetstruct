import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { apiService } from '../services/api';
import { User } from '../types';
import './SettingsPage.css';

export const SettingsPage: React.FC = () => {
  const { user, logout, updateUser } = useAuth();
  const navigate = useNavigate();
  const [name, setName] = useState(user?.Name || '');
  const [email, setEmail] = useState(user?.Email || '');
  const [password, setPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      if (!user) {
        setError('Пользователь не найден');
        setLoading(false);
        return;
      }

      const updatedUserData: User = {
        ...user,
        Name: name,
        Email: email,
      };
      
      if (newPassword) {
        if (!password) {
          setError('Текущий пароль требуется для смены пароля');
          setLoading(false);
          return;
        }
        updatedUserData.Password = newPassword;
      }

      const updatedUser = await apiService.updateUser(updatedUserData);
      updateUser(updatedUser);
      setSuccess('Профиль успешно обновлен');
      setPassword('');
      setNewPassword('');
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteAccount = async () => {
    setError('');
    setLoading(true);

    try {
      await apiService.deleteUser();
      logout();
      navigate('/login');
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
      setLoading(false);
    }
  };

  return (
    <div className="settings-page">
      <div className="settings-container">
        <h1>Настройки</h1>

        <div className="settings-section">
          <h2>Информация об аккаунте</h2>
          <form onSubmit={handleUpdateProfile} className="settings-form">
            {error && <div className="error-message">{error}</div>}
            {success && <div className="success-message">{success}</div>}

            <div className="form-group">
              <label>Имя</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label>Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label>Текущий пароль (требуется для смены пароля)</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Введите текущий пароль"
              />
            </div>

            <div className="form-group">
              <label>Новый пароль</label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="Введите новый пароль (оставьте пустым, чтобы оставить текущий)"
              />
            </div>

            <button type="submit" className="update-button" disabled={loading}>
              {loading ? 'Обновление...' : 'Обновить профиль'}
            </button>
          </form>
        </div>

        <div className="settings-section danger-zone">
          <h2>Опасная зона</h2>
          <div className="danger-actions">
            <div>
              <h3>Удалить аккаунт</h3>
              <p>Это действие нельзя отменить. Все ваши листы и данные будут безвозвратно удалены.</p>
            </div>
            {!showDeleteConfirm ? (
              <button
                className="delete-button"
                onClick={() => setShowDeleteConfirm(true)}
              >
                Удалить аккаунт
              </button>
            ) : (
              <div className="delete-confirm">
                <p>Вы уверены? Это действие нельзя отменить.</p>
                <div className="confirm-buttons">
                  <button
                    className="confirm-delete-button"
                    onClick={handleDeleteAccount}
                    disabled={loading}
                  >
                    Да, удалить
                  </button>
                  <button
                    className="cancel-button"
                    onClick={() => setShowDeleteConfirm(false)}
                  >
                    Отмена
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

