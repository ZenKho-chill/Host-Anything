import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { LogOut, Plus, Activity, Server, Box } from 'lucide-react';

interface Service {
  id: string;
  state: string;
}

export const DashboardView: React.FC = () => {
  const [services, setServices] = useState<Service[]>([]);
  const navigate = useNavigate();

  const fetchServices = async () => {
    try {
      const res = await fetch('/api/v1/services', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      if (res.status === 401) {
        localStorage.removeItem('token');
        navigate('/login');
        return;
      }
      const data = await res.json();
      setServices(data || []);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchServices();
    const interval = setInterval(fetchServices, 5000);
    return () => clearInterval(interval);
  }, [navigate]);

  return (
    <div className="w-full max-w-4xl mx-auto p-6 animate-fade-in" style={{ margin: '0 auto' }}>
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Server className="text-[var(--accent-primary)]" /> Dashboard
          </h1>
          <p className="text-[var(--text-muted)] mt-2">Manage your active deployments</p>
        </div>
        <div className="flex gap-4">
          <Button onClick={() => navigate('/templates')}>
            <Plus size={18} /> Deploy Service
          </Button>
          <Button variant="secondary" onClick={() => { localStorage.removeItem('token'); navigate('/login'); }}>
            <LogOut size={18} /> Logout
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4 mb-8">
        <Card glass className="flex flex-col gap-2">
          <div className="flex justify-between items-center text-[var(--text-secondary)]">
            <span>Total Services</span>
            <Box size={20} />
          </div>
          <span className="text-3xl font-bold">{services.length}</span>
        </Card>
        <Card glass className="flex flex-col gap-2">
          <div className="flex justify-between items-center text-[var(--text-secondary)]">
            <span>Running</span>
            <Activity size={20} className="text-[var(--success)]" />
          </div>
          <span className="text-3xl font-bold">
            {services.filter(s => s.state === 'RUNNING').length}
          </span>
        </Card>
      </div>

      <h2 className="text-2xl font-bold mb-4">Active Services</h2>
      {services.length === 0 ? (
        <Card glass className="flex flex-col items-center justify-center p-6 text-center" style={{ minHeight: '200px' }}>
          <Box size={48} className="text-[var(--text-muted)] mb-4" />
          <h3 className="text-xl font-medium mb-2">No services running</h3>
          <p className="text-[var(--text-muted)] mb-4">You haven't deployed anything yet.</p>
          <Button onClick={() => navigate('/templates')}>Browse Templates</Button>
        </Card>
      ) : (
        <div className="grid grid-cols-1 gap-4">
          {services.map(svc => (
            <Card key={svc.id} glass className="flex justify-between items-center">
              <div>
                <h3 className="font-bold text-lg">{svc.id}</h3>
                <div className="flex items-center text-sm text-[var(--text-secondary)] mt-2">
                  <span className={`status-dot ${svc.state.toLowerCase()}`}></span>
                  {svc.state}
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="secondary" size="sm">Logs</Button>
                <Button variant="danger" size="sm">Stop</Button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
};
