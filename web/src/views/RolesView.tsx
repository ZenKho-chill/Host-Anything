import React from 'react';
import { Card } from '../components/ui/Card';
import { Shield } from 'lucide-react';

export const RolesView: React.FC = () => {
  return (
    <div className="animate-fade-in">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Shield className="text-[var(--accent-primary)]" /> Roles & Permissions
          </h1>
          <p className="text-[var(--text-muted)] mt-2">Manage RBAC roles</p>
        </div>
      </div>
      <Card glass className="p-8 text-center text-[var(--text-muted)]">
        Role Management Module (Placeholder)
      </Card>
    </div>
  );
};
