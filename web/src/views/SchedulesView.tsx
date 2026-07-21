import React from 'react';
import { Card } from '../components/ui/Card';
import { Calendar } from 'lucide-react';

export const SchedulesView: React.FC = () => {
  return (
    <div className="animate-fade-in">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Calendar className="text-[var(--accent-primary)]" /> Schedules
          </h1>
          <p className="text-[var(--text-muted)] mt-2">Manage background cron jobs</p>
        </div>
      </div>
      <Card glass className="p-8 text-center text-[var(--text-muted)]">
        Schedules Management Module (Placeholder)
      </Card>
    </div>
  );
};
