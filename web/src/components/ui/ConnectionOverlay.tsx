import React, { useEffect, useState } from 'react';
import { WifiOff } from 'lucide-react';

export const ConnectionOverlay: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isOffline, setIsOffline] = useState(false);

  useEffect(() => {
    const originalFetch = window.fetch;

    window.fetch = async (...args) => {
      let [resource, config] = args;
      
      // Auto-inject Authorization header
      const token = localStorage.getItem('token');
      if (token) {
        if (typeof config === 'undefined') {
          config = {};
        }
        if (typeof config.headers === 'undefined') {
          config.headers = {};
        }
        
        // Convert headers to Headers object if it isn't already to easily append/set
        const headers = new Headers(config.headers);
        if (!headers.has('Authorization')) {
          headers.set('Authorization', `Bearer ${token}`);
        }
        
        config.headers = headers;
      }

      try {
        const response = await originalFetch(resource, config);
        
        // If the request succeeds, we are online
        if (isOffline) setIsOffline(false);

        // Handle 401 Unauthorized globally
        if (response.status === 401) {
          localStorage.removeItem('token');
          window.location.href = '/login';
        }

        return response;
      } catch (error) {
        // Network error / connection refused
        setIsOffline(true);
        throw error;
      }
    };

    return () => {
      window.fetch = originalFetch;
    };
  }, [isOffline]);

  useEffect(() => {
    if (!isOffline) return;

    const interval = setInterval(() => {
      // Use original fetch to avoid interceptor loop
      fetch('/api/v1/health')
        .then(res => {
          if (res.ok) {
            setIsOffline(false);
            window.location.reload(); // Reload to get fresh state
          }
        })
        .catch(() => {
          // Still offline
        });
    }, 3000);

    return () => clearInterval(interval);
  }, [isOffline]);

  return (
    <>
      {children}
      {isOffline && (
        <div style={{
          position: 'fixed',
          top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(18, 20, 28, 0.8)',
          backdropFilter: 'blur(8px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
          color: 'var(--error)'
        }}>
          <div style={{
            background: 'var(--bg-secondary)',
            padding: '32px',
            borderRadius: '16px',
            border: '1px solid rgba(239, 68, 68, 0.2)',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '16px',
            boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.5)'
          }}>
            <div className="animate-pulse">
              <WifiOff size={48} />
            </div>
            <h2 className="text-2xl font-bold text-[var(--text-primary)]">Connection Lost</h2>
            <p className="text-[var(--text-muted)] text-center max-w-sm">
              The Host Anything core daemon is currently unreachable. <br/>
              Attempting to reconnect...
            </p>
          </div>
        </div>
      )}
    </>
  );
};
