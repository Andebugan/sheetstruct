import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Navigation } from './components/Navigation';
import { LandingPage } from './pages/LandingPage';
import { SettingsPage } from './pages/SettingsPage';
import { SheetManagementPage } from './pages/SheetManagementPage';
import { SheetBuildingPage } from './pages/SheetBuildingPage';
import { ResetPasswordPage } from './pages/ResetPasswordPage';
import { useAuth } from './context/AuthContext';
import './App.css';

const AppRoutes: React.FC = () => {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div className="loading-screen">Loading...</div>;
  }

  return (
    <Routes>
      <Route
        path="/login"
        element={isAuthenticated ? <Navigate to="/sheets" replace /> : <LandingPage />}
      />
      <Route
        path="/reset-password"
        element={<ResetPasswordPage />}
      />
      <Route
        path="/sheets"
        element={
          <ProtectedRoute>
            <Navigation />
            <SheetManagementPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/sheet/:id"
        element={
          <ProtectedRoute>
            <SheetBuildingPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/settings"
        element={
          <ProtectedRoute>
            <Navigation />
            <SettingsPage />
          </ProtectedRoute>
        }
      />
      <Route path="/" element={<Navigate to="/sheets" replace />} />
    </Routes>
  );
};

function App() {
  return (
    <Router>
      <AuthProvider>
        <div className="App">
          <AppRoutes />
        </div>
      </AuthProvider>
    </Router>
  );
}

export default App;
