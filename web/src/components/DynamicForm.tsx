import React, { useEffect, useState } from 'react';
import { Button } from './ui/Button';
import { Input } from './ui/Input';

interface ConfigVar {
  description: string;
  type: string;
  default?: any;
  secret?: boolean;
}

interface TemplateDetail {
  meta: any;
  config: Record<string, ConfigVar>;
}

export const DynamicForm: React.FC<{ templateName: string; onSuccess: () => void }> = ({ templateName, onSuccess }) => {
  const [template, setTemplate] = useState<TemplateDetail | null>(null);
  const [formData, setFormData] = useState<Record<string, string>>({});
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    fetch(`/api/v1/templates/${templateName}`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
      .then(res => res.json())
      .then(data => {
        setTemplate(data);
        // Initialize defaults
        const defaults: Record<string, string> = {};
        if (data.config) {
          Object.entries(data.config).forEach(([key, val]: [string, any]) => {
            if (val.default !== undefined) {
              defaults[key] = String(val.default);
            }
          });
        }
        setFormData(defaults);
      })
      .catch(console.error);
  }, [templateName]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const res = await fetch('/api/v1/services', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          template_name: templateName,
          config: formData
        })
      });

      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || 'Deployment failed');
      }

      onSuccess();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (!template) {
    return <div className="text-center p-4">Loading template schema...</div>;
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      {template.config && Object.entries(template.config).map(([key, field]) => (
        <Input
          key={key}
          label={key}
          type={field.secret ? 'password' : (field.type === 'int' ? 'number' : 'text')}
          helpText={field.description}
          value={formData[key] || ''}
          onChange={(e) => setFormData({ ...formData, [key]: e.target.value })}
          required={field.default === undefined}
        />
      ))}

      {error && (
        <div className="p-3 rounded bg-[var(--error-bg)] border border-[var(--error)]">
          <p className="text-[var(--error)] text-sm font-medium">{error}</p>
        </div>
      )}

      <div className="flex justify-end gap-2 mt-4">
        <Button type="submit" isLoading={isLoading}>
          Deploy Service
        </Button>
      </div>
    </form>
  );
};
