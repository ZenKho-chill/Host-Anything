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
          backgroundColor: 'rgba(10, 12, 16, 0.85)',
          backdropFilter: 'blur(16px)',
          WebkitBackdropFilter: 'blur(16px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
          color: 'var(--error)'
        }}>
          {/* Pulsing Background Glow */}
          <div style={{
            position: 'absolute',
            width: '300px', height: '300px',
            background: 'radial-gradient(circle, rgba(239, 68, 68, 0.15) 0%, transparent 70%)',
            borderRadius: '50%',
            animation: 'pulse-subtle 4s infinite'
          }}></div>

          <div className="animate-slide-up animate-float" style={{
            background: 'linear-gradient(145deg, rgba(30, 33, 43, 0.9), rgba(20, 22, 30, 0.9))',
            padding: '48px',
            borderRadius: '24px',
            border: '1px solid rgba(239, 68, 68, 0.3)',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '24px',
            boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.7), 0 0 0 1px rgba(255,255,255,0.05) inset',
            position: 'relative',
            overflow: 'hidden'
          }}>
            
            {/* Radar / Ping Effect */}
            <div style={{ position: 'relative', width: '80px', height: '80px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <div className="animate-ping-large" style={{
                position: 'absolute',
                width: '100%', height: '100%',
                borderRadius: '50%',
                backgroundColor: 'rgba(239, 68, 68, 0.4)',
              }}></div>
              <div style={{
                position: 'absolute',
                width: '100%', height: '100%',
                borderRadius: '50%',
                backgroundColor: 'rgba(239, 68, 68, 0.1)',
                border: '2px solid rgba(239, 68, 68, 0.5)'
              }}></div>
              <WifiOff size={40} style={{ color: '#fca5a5', zIndex: 2 }} />
            </div>

            <div style={{ textAlign: 'center', zIndex: 2 }}>
              <h2 style={{ 
                fontSize: '28px', 
                fontWeight: '800', 
                color: '#fff', 
                marginBottom: '8px',
                letterSpacing: '-0.05em',
                textShadow: '0 2px 10px rgba(239,68,68,0.5)'
              }}>
                SYSTEM OFFLINE
              </h2>
              <p style={{ 
                color: 'var(--text-secondary)', 
                fontSize: '15px', 
                maxWidth: '320px', 
                lineHeight: '1.6' 
              }}>
                The Host Anything core daemon is unreachable. We are actively attempting to restore the connection.
              </p>
            </div>
            
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '12px',
              padding: '12px 24px',
              background: 'rgba(0,0,0,0.3)',
              borderRadius: '99px',
              border: '1px solid rgba(255,255,255,0.05)',
              marginTop: '8px'
            }}>
              <div className="status-dot deploying"></div>
              <span style={{ color: 'var(--text-secondary)', fontSize: '14px', fontWeight: '500', letterSpacing: '0.05em' }}>RECONNECTING...</span>
            </div>

          </div>
        </div>
      )}
    </>
  );
};
