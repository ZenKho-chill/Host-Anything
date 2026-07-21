import React from 'react';
import { Card } from '../components/ui/Card';
import { Users, Plus } from 'lucide-react';
import { Button } from '../components/ui/Button';

export const UsersView: React.FC = () => {
  return (
    <div className="animate-fade-in">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Users className="text-[var(--accent-primary)]" /> Users Management
          </h1>
          <p className="text-[var(--text-muted)] mt-2">Manage users and their roles</p>
        </div>
        <Button><Plus size={18} /> Add User</Button>
      </div>
      <Card glass className="p-8 text-center text-[var(--text-muted)]">
        User Management Module (Placeholder)
      </Card>
    </div>
  );
};
