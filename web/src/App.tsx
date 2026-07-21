import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { LoginView } from './views/LoginView';
import { DashboardView } from './views/DashboardView';
import { TemplateBrowserView } from './views/TemplateBrowserView';
import { MarketplaceView } from './views/MarketplaceView';
import { UsersView } from './views/UsersView';
import { RolesView } from './views/RolesView';
import { SchedulesView } from './views/SchedulesView';
import { SettingsView } from './views/SettingsView';
import { FilesView } from './views/FilesView';
import { Layout } from './components/ui/Layout';

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
            <Layout><DashboardView /></Layout>
          </PrivateRoute>
        } />
        
        <Route path="/templates" element={
          <PrivateRoute>
            <Layout><TemplateBrowserView /></Layout>
          </PrivateRoute>
        } />

        <Route path="/marketplace" element={
          <PrivateRoute>
            <Layout><MarketplaceView /></Layout>
          </PrivateRoute>
        } />

        <Route path="/users" element={
          <PrivateRoute>
            <Layout><UsersView /></Layout>
          </PrivateRoute>
        } />

        <Route path="/roles" element={
          <PrivateRoute>
            <Layout><RolesView /></Layout>
          </PrivateRoute>
        } />

        <Route path="/schedules" element={
          <PrivateRoute>
            <Layout><SchedulesView /></Layout>
          </PrivateRoute>
        } />

        <Route path="/settings" element={
          <PrivateRoute>
            <Layout><SettingsView /></Layout>
          </PrivateRoute>
        } />

        <Route path="/files" element={
          <PrivateRoute>
            <Layout><FilesView /></Layout>
          </PrivateRoute>
        } />
        
        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;
