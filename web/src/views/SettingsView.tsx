import React, { useEffect, useState } from 'react';
import { Card } from '../components/ui/Card';
import { Settings, Save } from 'lucide-react';
import { Button } from '../components/ui/Button';

export const SettingsView: React.FC = () => {
  const [timezone, setTimezone] = useState('');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetch('/api/v1/settings', { headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }})
      .then(r => r.json())
      .then(d => setTimezone(d.timezone || 'UTC'));
  }, []);

  const saveSettings = async () => {
    setSaving(true);
    await fetch('/api/v1/settings', {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json' 
      },
      body: JSON.stringify({ timezone })
    });
    setSaving(false);
  };

  return (
    <div className="animate-fade-in">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Settings className="text-[var(--accent-primary)]" /> System Settings
          </h1>
          <p className="text-[var(--text-muted)] mt-2">Configure global daemon settings</p>
        </div>
      </div>
      <Card glass className="max-w-2xl">
        <h3 className="text-lg font-bold mb-4">Localization</h3>
        <div className="mb-4">
          <label className="block text-sm font-medium mb-1">Timezone</label>
          <select 
            value={timezone} 
            onChange={e => setTimezone(e.target.value)}
            style={{ 
              width: '100%', padding: '10px', 
              background: 'var(--surface-elevated)', border: '1px solid var(--border)',
              borderRadius: '8px', color: 'var(--text-primary)'
            }}
          >
            <option value="UTC">UTC</option>
            <option value="America/New_York">America/New_York</option>
            <option value="Europe/London">Europe/London</option>
            <option value="Asia/Ho_Chi_Minh">Asia/Ho_Chi_Minh</option>
          </select>
        </div>
        <Button onClick={saveSettings} disabled={saving}>
          <Save size={18} /> {saving ? 'Saving...' : 'Save Settings'}
        </Button>
      </Card>
    </div>
  );
};
