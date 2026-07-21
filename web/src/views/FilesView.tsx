import React, { useEffect, useState } from 'react';
import { Card } from '../components/ui/Card';
import { Folder, File, ArrowLeft, Download, Home } from 'lucide-react';
import { Button } from '../components/ui/Button';

interface FileInfo {
  name: string;
  is_dir: boolean;
  size: number;
}

export const FilesView: React.FC = () => {
  const [path, setPath] = useState('/');
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [error, setError] = useState('');

  const fetchFiles = async (p: string) => {
    try {
      const url = `/api/v1/files${p.startsWith('/') ? p : '/' + p}`;
      const res = await fetch(url, { headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }});
      if (res.ok) {
        const data = await res.json();
        setFiles(data || []);
        setError('');
      } else {
        const err = await res.json();
        setError(err.error || 'Access Denied');
      }
    } catch {
      setError('Network error');
    }
  };

  useEffect(() => {
    fetchFiles(path);
  }, [path]);

  const goUp = () => {
    const parts = path.split('/').filter(Boolean);
    parts.pop();
    setPath('/' + parts.join('/'));
  };

  const handleOpen = (f: FileInfo) => {
    if (f.is_dir) {
      setPath(path === '/' ? `/${f.name}` : `${path}/${f.name}`);
    } else {
      // Open file in new tab (would require a temporary token in a real app, 
      // but for this milestone we just alert)
      alert('File download requires token integration in URL.');
    }
  };

  return (
    <div className="animate-fade-in">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Folder className="text-[var(--accent-primary)]" /> File Manager
          </h1>
          <p className="text-[var(--text-muted)] mt-2">Manage service data and volumes</p>
        </div>
      </div>

      <Card glass>
        <div className="flex items-center gap-4 mb-4 pb-4 border-b border-[var(--border)]">
          <Button variant="secondary" onClick={() => setPath('/')} disabled={path === '/'}><Home size={18} /></Button>
          <Button variant="secondary" onClick={goUp} disabled={path === '/'}><ArrowLeft size={18} /> Up</Button>
          <div className="font-mono bg-[var(--surface-elevated)] px-4 py-2 rounded flex-1">
            {path}
          </div>
        </div>

        {error ? (
          <div className="p-8 text-center text-[var(--error)]">{error}</div>
        ) : (
          <div className="grid grid-cols-1 gap-2">
            {files.map(f => (
              <div 
                key={f.name} 
                onClick={() => handleOpen(f)}
                className="flex items-center justify-between p-3 rounded hover:bg-[var(--surface-elevated)] cursor-pointer transition-colors"
              >
                <div className="flex items-center gap-3">
                  {f.is_dir ? <Folder className="text-yellow-400" /> : <File className="text-[var(--text-muted)]" />}
                  <span className="font-medium">{f.name}</span>
                </div>
                <div className="flex items-center gap-4 text-sm text-[var(--text-muted)]">
                  {!f.is_dir && <span>{(f.size / 1024).toFixed(2)} KB</span>}
                  {!f.is_dir && <Button variant="secondary" size="sm"><Download size={14} /></Button>}
                </div>
              </div>
            ))}
            {files.length === 0 && (
              <div className="p-8 text-center text-[var(--text-muted)]">Directory is empty</div>
            )}
          </div>
        )}
      </Card>
    </div>
  );
};
