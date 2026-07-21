import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card } from '../components/ui/Card';
import { Shield, AlertCircle } from 'lucide-react';
import { setToken, getToken } from '../utils/auth';

export const LoginView: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);
  const [error, setError] = useState('');
  const [isShaking, setIsShaking] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  // Auto redirect if already logged in
  useEffect(() => {
    if (getToken()) {
      navigate('/');
    }
  }, [navigate]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setIsShaking(false);

    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || 'Invalid credentials');
      }

      setToken(data.token, rememberMe);
      navigate('/');
    } catch (err: any) {
      setError(err.message);
      setIsShaking(true);
      // Remove shake class after animation completes so it can trigger again
      setTimeout(() => setIsShaking(false), 500);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-container" style={{
      background: 'radial-gradient(circle at center, var(--bg-secondary) 0%, var(--bg-primary) 100%)'
    }}>
      <div className="w-full max-w-md p-6 relative">
        {/* Glow Effects behind the card */}
        <div style={{
          position: 'absolute',
          top: '20%', left: '10%',
          width: '200px', height: '200px',
          background: 'var(--accent-primary)',
          filter: 'blur(100px)',
          opacity: 0.15,
          zIndex: 0
        }}></div>

        <div className="flex flex-col items-center mb-8 animate-fade-in relative z-10">
          <div style={{
            background: 'linear-gradient(135deg, var(--accent-primary), var(--accent-hover))',
            padding: '16px',
            borderRadius: '20px',
            marginBottom: '20px',
            boxShadow: '0 10px 30px -10px var(--accent-primary)'
          }}>
            <Shield className="text-white" size={40} />
          </div>
          <h1 className="text-4xl font-extrabold tracking-tight" style={{ 
            background: 'linear-gradient(to right, #fff, var(--text-secondary))',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent'
          }}>
            Host Anything
          </h1>
          <p className="text-[var(--text-muted)] mt-3 font-medium tracking-wide">SECURE ACCESS TERMINAL</p>
        </div>

        <Card glass className={`relative z-10 transition-transform ${isShaking ? 'animate-shake' : 'animate-slide-up'}`} 
              style={{ border: error ? '1px solid var(--error)' : '1px solid rgba(255,255,255,0.05)' }}>
          <form onSubmit={handleLogin} className="flex flex-col gap-5 p-2">
            <Input
              label="Username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              autoFocus
              className={error ? 'input-error' : ''}
            />
            <Input
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className={error ? 'input-error' : ''}
            />
            
            <div className="flex items-center justify-between mt-1 mb-2">
              <label className="flex items-center gap-2 cursor-pointer group">
                <div style={{
                  width: '18px', height: '18px',
                  borderRadius: '4px',
                  border: rememberMe ? 'none' : '1px solid var(--border)',
                  background: rememberMe ? 'var(--accent-primary)' : 'var(--bg-primary)',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  transition: 'all 0.2s ease'
                }}>
                  {rememberMe && <svg viewBox="0 0 14 14" fill="none" xmlns="http://www.w3.org/2000/svg" style={{ width: '12px', height: '12px' }}>
                    <path d="M11.6666 3.5L5.24992 9.91667L2.33325 7" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>}
                </div>
                <input 
                  type="checkbox" 
                  checked={rememberMe} 
                  onChange={(e) => setRememberMe(e.target.checked)} 
                  className="hidden"
                />
                <span className="text-sm font-medium text-[var(--text-secondary)] group-hover:text-[var(--text-primary)] transition-colors">
                  Remember me
                </span>
              </label>
            </div>
            
            {error && (
              <div className="p-3 rounded-lg flex items-center gap-3 animate-fade-in" style={{
                background: 'rgba(239, 68, 68, 0.1)',
                borderLeft: '4px solid var(--error)'
              }}>
                <AlertCircle size={20} className="text-[var(--error)] flex-shrink-0" />
                <p className="text-[var(--error)] text-sm font-medium">{error}</p>
              </div>
            )}
            
            <Button type="submit" size="lg" isLoading={isLoading} className="mt-2 w-full font-bold tracking-wide"
                    style={{ background: 'var(--accent-primary)', boxShadow: '0 4px 14px 0 rgba(99, 102, 241, 0.39)' }}>
              AUTHENTICATE
            </Button>
          </form>
        </Card>
      </div>
    </div>
  );
};
