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
              background: 'var(--bg-elevated)', border: '1px solid var(--border)',
              borderRadius: '8px', color: 'var(--text-primary)'
            }}
          >
            <optgroup label="UTC Offsets">
              <option value="UTC-12">UTC-12</option>
              <option value="UTC-11">UTC-11</option>
              <option value="UTC-10">UTC-10</option>
              <option value="UTC-09">UTC-09</option>
              <option value="UTC-08">UTC-08</option>
              <option value="UTC-07">UTC-07</option>
              <option value="UTC-06">UTC-06</option>
              <option value="UTC-05">UTC-05</option>
              <option value="UTC-04">UTC-04</option>
              <option value="UTC-03">UTC-03</option>
              <option value="UTC-02">UTC-02</option>
              <option value="UTC-01">UTC-01</option>
              <option value="UTC">UTC (00)</option>
              <option value="UTC+01">UTC+01</option>
              <option value="UTC+02">UTC+02</option>
              <option value="UTC+03">UTC+03</option>
              <option value="UTC+04">UTC+04</option>
              <option value="UTC+05">UTC+05</option>
              <option value="UTC+06">UTC+06</option>
              <option value="UTC+07">UTC+07 (Vietnam, Indochina)</option>
              <option value="UTC+08">UTC+08 (China, Singapore)</option>
              <option value="UTC+09">UTC+09 (Japan, Korea)</option>
              <option value="UTC+10">UTC+10</option>
              <option value="UTC+11">UTC+11</option>
              <option value="UTC+12">UTC+12</option>
              <option value="UTC+13">UTC+13</option>
              <option value="UTC+14">UTC+14</option>
            </optgroup>
            <optgroup label="Common Cities">
              <option value="America/Los_Angeles">America/Los Angeles</option>
              <option value="America/New_York">America/New York</option>
              <option value="Europe/London">Europe/London</option>
              <option value="Europe/Paris">Europe/Paris</option>
              <option value="Asia/Ho_Chi_Minh">Asia/Ho Chi Minh</option>
              <option value="Asia/Tokyo">Asia/Tokyo</option>
              <option value="Australia/Sydney">Australia/Sydney</option>
            </optgroup>
          </select>
        </div>
        <Button onClick={saveSettings} disabled={saving}>
          <Save size={18} /> {saving ? 'Saving...' : 'Save Settings'}
        </Button>
      </Card>
    </div>
  );
};
