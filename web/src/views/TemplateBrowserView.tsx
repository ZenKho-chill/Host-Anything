import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { ArrowLeft, Box, Download } from 'lucide-react';
import { DynamicForm } from '../components/DynamicForm';
import { Modal } from '../components/ui/Modal';

interface TemplateSummary {
  name: string;
  version: string;
  description: string;
  author: string;
  tags: string[];
}

export const TemplateBrowserView: React.FC = () => {
  const [templates, setTemplates] = useState<TemplateSummary[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetch('/api/v1/templates', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
      .then(res => res.json())
      .then(data => setTemplates(data || []))
      .catch(console.error);
  }, []);

  return (
    <div className="w-full max-w-4xl mx-auto p-6 animate-fade-in" style={{ margin: '0 auto' }}>
      <div className="flex items-center gap-4 mb-8">
        <Button variant="ghost" onClick={() => navigate('/')} className="px-2">
          <ArrowLeft size={20} />
        </Button>
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Box className="text-[var(--accent-primary)]" /> Template Browser
          </h1>
          <p className="text-[var(--text-muted)] mt-2">Select a template to deploy a new service</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {templates.map(tmpl => (
          <Card key={tmpl.name} glass className="flex flex-col">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h3 className="font-bold text-xl">{tmpl.name}</h3>
                <span className="text-sm text-[var(--text-muted)]">v{tmpl.version} by {tmpl.author}</span>
              </div>
            </div>
            <p className="text-[var(--text-secondary)] mb-6 flex-grow">{tmpl.description}</p>
            <div className="flex gap-2 flex-wrap mb-4">
              {tmpl.tags?.map(tag => (
                <span key={tag} className="text-xs bg-[var(--bg-elevated)] px-2 py-1 rounded border border-[var(--border)]">
                  {tag}
                </span>
              ))}
            </div>
            <Button className="w-full" onClick={() => setSelectedTemplate(tmpl.name)}>
              <Download size={18} /> Configure & Deploy
            </Button>
          </Card>
        ))}
      </div>

      <Modal 
        isOpen={!!selectedTemplate} 
        onClose={() => setSelectedTemplate(null)}
        title={`Deploy ${selectedTemplate}`}
      >
        {selectedTemplate && (
          <DynamicForm 
            templateName={selectedTemplate} 
            onSuccess={() => navigate('/')} 
          />
        )}
      </Modal>
    </div>
  );
};
