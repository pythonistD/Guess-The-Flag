import React from 'react';
import { Navigate } from 'react-router-dom';
import { TokenUtils } from '../services/api';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const isAuthenticated = TokenUtils.isAuthenticated();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />; // relative to Router basename
  }

  return <>{children}</>;
};

export default ProtectedRoute; 