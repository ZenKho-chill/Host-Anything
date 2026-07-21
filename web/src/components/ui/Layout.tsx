import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Server, Users, Shield, Calendar, Settings, Folder, LayoutDashboard, ShoppingBag, LogOut } from 'lucide-react';

interface LayoutProps {
  children: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const navItems = [
    { name: 'Dashboard', path: '/', icon: <LayoutDashboard size={20} /> },
    { name: 'Marketplace', path: '/marketplace', icon: <ShoppingBag size={20} /> },
    { name: 'File Manager', path: '/files', icon: <Folder size={20} /> },
    { name: 'Users', path: '/users', icon: <Users size={20} /> },
    { name: 'Roles', path: '/roles', icon: <Shield size={20} /> },
    { name: 'Schedules', path: '/schedules', icon: <Calendar size={20} /> },
    { name: 'Settings', path: '/settings', icon: <Settings size={20} /> },
  ];

  return (
    <div style={{ display: 'flex', minHeight: '100vh', background: 'var(--bg)' }}>
      {/* Sidebar */}
      <div style={{ 
        width: '260px', 
        background: 'var(--surface)', 
        borderRight: '1px solid var(--border)',
        display: 'flex',
        flexDirection: 'column'
      }}>
        <div style={{ padding: '24px', display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div style={{ 
            width: '40px', height: '40px', 
            borderRadius: '10px', 
            background: 'linear-gradient(135deg, var(--accent-primary), var(--accent-hover))',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            color: '#fff'
          }}>
            <Server size={24} />
          </div>
          <h2 style={{ fontSize: '18px', fontWeight: 'bold', margin: 0 }}>Host Anything</h2>
        </div>

        <nav style={{ flex: 1, padding: '0 12px' }}>
          {navItems.map(item => {
            const active = location.pathname === item.path || (item.path !== '/' && location.pathname.startsWith(item.path));
            return (
              <button
                key={item.path}
                onClick={() => navigate(item.path)}
                style={{
                  width: '100%',
                  display: 'flex', alignItems: 'center', gap: '12px',
                  padding: '12px 16px',
                  marginBottom: '4px',
                  borderRadius: '8px',
                  background: active ? 'var(--surface-elevated)' : 'transparent',
                  color: active ? 'var(--text-primary)' : 'var(--text-muted)',
                  border: 'none',
                  cursor: 'pointer',
                  textAlign: 'left',
                  fontWeight: active ? 600 : 500,
                  transition: 'all 0.2s ease'
                }}
              >
                {item.icon}
                {item.name}
              </button>
            );
          })}
        </nav>

        <div style={{ padding: '24px 12px' }}>
          <button
            onClick={() => { localStorage.removeItem('token'); navigate('/login'); }}
            style={{
              width: '100%',
              display: 'flex', alignItems: 'center', gap: '12px',
              padding: '12px 16px',
              borderRadius: '8px',
              background: 'transparent',
              color: 'var(--error)',
              border: 'none',
              cursor: 'pointer',
              textAlign: 'left',
              fontWeight: 500,
              transition: 'background 0.2s ease'
            }}
          >
            <LogOut size={20} />
            Logout
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div style={{ flex: 1, padding: '32px', overflowY: 'auto' }}>
        {children}
      </div>
    </div>
  );
};
