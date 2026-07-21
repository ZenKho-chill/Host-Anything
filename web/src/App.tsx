import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { LoginView } from './views/LoginView';
import { DashboardView } from './views/DashboardView';
import { TemplateBrowserView } from './views/TemplateBrowserView';
import { MarketplaceView } from './views/MarketplaceView';

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const token = localStorage.getItem('token');
  if (!token) {
    return <Navigate to="/login" />;
  }
  return <>{children}</>;
};

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginView />} />
        
        {/* Protected Routes */}
        <Route path="/" element={
          <PrivateRoute>
            <DashboardView />
          </PrivateRoute>
        } />
        
        <Route path="/templates" element={
          <PrivateRoute>
            <TemplateBrowserView />
          </PrivateRoute>
        } />

        <Route path="/marketplace" element={
          <PrivateRoute>
            <MarketplaceView />
          </PrivateRoute>
        } />
        
        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;
