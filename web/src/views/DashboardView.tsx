import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { Plus, Activity, Server, Box, ShoppingBag, Cpu, HardDrive, Clock } from 'lucide-react';

interface Service {
  id: string;
  state: string;
}

export const DashboardView: React.FC = () => {
  const [services, setServices] = useState<Service[]>([]);
  const navigate = useNavigate();

  const fetchServices = async () => {
    try {
      // The interceptor in ConnectionOverlay handles the token injection and 401s
      const res = await fetch('/api/v1/services');
      if (res.ok) {
        const data = await res.json();
        setServices(data || []);
      }
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
    <div className="animate-fade-in" style={{ width: '100%', maxWidth: '1400px', margin: '0 auto' }}>
      {/* Header Area */}
      <div className="flex justify-between items-end mb-8 border-b border-[var(--border)] pb-6">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div style={{
              background: 'linear-gradient(135deg, var(--accent-primary), var(--accent-hover))',
              padding: '10px',
              borderRadius: '12px',
              boxShadow: '0 4px 15px -3px var(--accent-primary)'
            }}>
              <Server className="text-white" size={24} />
            </div>
            <h1 className="text-4xl font-extrabold tracking-tight">System Overview</h1>
          </div>
          <p className="text-[var(--text-muted)] text-lg">Real-time status and telemetry of your deployments</p>
        </div>
        <div className="flex gap-3">
          <Button variant="secondary" onClick={() => navigate('/marketplace')} style={{ border: '1px solid var(--border)' }}>
            <ShoppingBag size={18} /> Marketplace
          </Button>
          <Button onClick={() => navigate('/templates')} style={{ background: 'var(--accent-primary)', boxShadow: '0 4px 14px 0 rgba(99, 102, 241, 0.3)' }}>
            <Plus size={18} /> Deploy Instance
          </Button>
        </div>
      </div>

      {/* Top Telemetry / Stats Grid */}
      <div className="grid grid-cols-4 gap-6 mb-10">
        <Card glass style={{ background: 'linear-gradient(145deg, var(--bg-secondary), var(--bg-primary))', borderTop: '2px solid var(--accent-primary)' }}>
          <div className="flex justify-between items-start mb-4 text-[var(--text-secondary)]">
            <span className="font-semibold text-sm uppercase tracking-wider">Total Services</span>
            <div className="p-2 bg-[var(--bg-elevated)] rounded-lg"><Box size={18} className="text-[var(--accent-primary)]" /></div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-4xl font-black">{services.length}</span>
            <span className="text-sm text-[var(--text-muted)]">deployed</span>
          </div>
        </Card>
        
        <Card glass style={{ background: 'linear-gradient(145deg, var(--bg-secondary), var(--bg-primary))', borderTop: '2px solid var(--success)' }}>
          <div className="flex justify-between items-start mb-4 text-[var(--text-secondary)]">
            <span className="font-semibold text-sm uppercase tracking-wider">Running</span>
            <div className="p-2 bg-[var(--bg-elevated)] rounded-lg"><Activity size={18} className="text-[var(--success)]" /></div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-4xl font-black">{services.filter(s => s.state === 'RUNNING').length}</span>
            <span className="text-sm text-[var(--text-muted)]">online</span>
          </div>
        </Card>

        <Card glass style={{ background: 'linear-gradient(145deg, var(--bg-secondary), var(--bg-primary))', borderTop: '2px solid var(--warning)' }}>
          <div className="flex justify-between items-start mb-4 text-[var(--text-secondary)]">
            <span className="font-semibold text-sm uppercase tracking-wider">CPU Load</span>
            <div className="p-2 bg-[var(--bg-elevated)] rounded-lg"><Cpu size={18} className="text-[var(--warning)]" /></div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-4xl font-black">12%</span>
            <span className="text-sm text-[var(--text-muted)]">avg</span>
          </div>
        </Card>

        <Card glass style={{ background: 'linear-gradient(145deg, var(--bg-secondary), var(--bg-primary))', borderTop: '2px solid #8b5cf6' }}>
          <div className="flex justify-between items-start mb-4 text-[var(--text-secondary)]">
            <span className="font-semibold text-sm uppercase tracking-wider">Uptime</span>
            <div className="p-2 bg-[var(--bg-elevated)] rounded-lg"><Clock size={18} className="text-[#8b5cf6]" /></div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-4xl font-black">99.9%</span>
            <span className="text-sm text-[var(--text-muted)]">SLA</span>
          </div>
        </Card>
      </div>

      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold tracking-tight">Active Deployments</h2>
      </div>

      {services.length === 0 ? (
        <Card glass className="flex flex-col items-center justify-center p-12 text-center border-dashed border-2 border-[var(--border)] bg-transparent">
          <div className="w-20 h-20 bg-[var(--bg-elevated)] rounded-full flex items-center justify-center mb-6">
            <Box size={40} className="text-[var(--text-muted)]" />
          </div>
          <h3 className="text-2xl font-bold mb-3">No instances running</h3>
          <p className="text-[var(--text-muted)] mb-8 max-w-md text-lg">Your workspace is empty. Deploy a new instance from a template or the marketplace to get started.</p>
          <Button onClick={() => navigate('/templates')} size="lg" style={{ background: 'var(--accent-primary)' }}>
            <Plus size={20} className="mr-2" /> Start First Deployment
          </Button>
        </Card>
      ) : (
        <div className="grid grid-cols-2 gap-6">
          {services.map(svc => (
            <Card key={svc.id} glass className="flex flex-col transition-all hover:-translate-y-1" style={{
              background: 'linear-gradient(180deg, var(--bg-secondary) 0%, rgba(30, 33, 43, 0.4) 100%)',
              border: '1px solid rgba(255,255,255,0.05)',
              boxShadow: '0 10px 25px -5px rgba(0,0,0,0.3)'
            }}>
              <div className="flex justify-between items-start mb-6">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-xl bg-[var(--bg-elevated)] flex items-center justify-center border border-[var(--border)]">
                    <HardDrive size={24} className="text-[var(--text-secondary)]" />
                  </div>
                  <div>
                    <h3 className="font-bold text-xl tracking-tight mb-1">{svc.id}</h3>
                    <div className="flex items-center gap-2">
                      <span className={`status-dot ${svc.state.toLowerCase()}`}></span>
                      <span className="text-sm font-medium text-[var(--text-secondary)] uppercase tracking-wider">{svc.state}</span>
                    </div>
                  </div>
                </div>
                <div className="bg-[var(--bg-elevated)] px-3 py-1 rounded-full text-xs font-mono text-[var(--text-muted)]">
                  ID: {svc.id.substring(0, 8)}
                </div>
              </div>
              
              <div className="mt-auto pt-4 border-t border-[var(--border)] flex gap-3 justify-end">
                <Button variant="ghost" size="sm" className="font-medium">Configure</Button>
                <Button variant="secondary" size="sm" className="font-medium border-[var(--border)]">Console Logs</Button>
                <Button variant="danger" size="sm" className="font-medium bg-[rgba(239,68,68,0.1)] text-[var(--error)] hover:bg-[var(--error)] hover:text-white border-transparent">Stop Instance</Button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
};
