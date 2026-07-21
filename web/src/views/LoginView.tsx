import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card } from '../components/ui/Card';
import { Shield } from 'lucide-react';

export const LoginView: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || 'Login failed');
      }

      localStorage.setItem('token', data.token);
      navigate('/');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="w-full max-w-md p-6">
        <div className="flex flex-col items-center mb-8 animate-fade-in">
          <div className="bg-[var(--accent-light)] p-3 rounded-full mb-4">
            <Shield className="text-[var(--accent-primary)]" size={32} />
          </div>
          <h1 className="text-3xl font-bold">Host Anything</h1>
          <p className="text-[var(--text-muted)] mt-2">Deploy anything. Configure everything.</p>
        </div>

        <Card glass className="animate-slide-up">
          <form onSubmit={handleLogin} className="flex flex-col gap-4">
            <Input
              label="Username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              autoFocus
            />
            <Input
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            
            {error && (
              <div className="p-3 rounded bg-[var(--error-bg)] border border-[var(--error)]">
                <p className="text-[var(--error)] text-sm font-medium">{error}</p>
              </div>
            )}
            
            <Button type="submit" size="lg" isLoading={isLoading} className="mt-2 w-full">
              Sign In
            </Button>
          </form>
        </Card>
      </div>
    </div>
  );
};
