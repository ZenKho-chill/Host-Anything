import React, { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { ArrowLeft, Search, Star, ShieldCheck, Users, ExternalLink, Download, X, ChevronRight } from 'lucide-react';

interface MarketplaceResult {
  name: string;
  owner: string;
  description: string;
  stars: number;
  repo_url: string;
  is_official: boolean;
}

interface TemplatePreview {
  template: {
    meta: { name: string; version: string; description: string; author: string; tags: string[] };
    runtime: { image: string; preferred: string; supported: string[] };
    config: Array<{ name: string; type: string; description: string; required: boolean; default?: unknown }>;
  };
  is_official: boolean;
}

const authHeader = () => ({ 'Authorization': `Bearer ${localStorage.getItem('token')}` });

export const MarketplaceView: React.FC = () => {
  const navigate = useNavigate();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<MarketplaceResult[]>([]);
  const [searching, setSearching] = useState(false);
  const [searched, setSearched] = useState(false);
  const [preview, setPreview] = useState<TemplatePreview | null>(null);
  const [selectedResult, setSelectedResult] = useState<MarketplaceResult | null>(null);
  const [loadingPreview, setLoadingPreview] = useState(false);
  const [installing, setInstalling] = useState(false);
  const [installMsg, setInstallMsg] = useState('');

  const search = useCallback(async () => {
    if (!query.trim()) return;
    setSearching(true);
    setSearched(true);
    try {
      const res = await fetch(`/api/v1/marketplace/search?q=${encodeURIComponent(query)}`, {
        headers: authHeader(),
      });
      if (res.status === 401) { localStorage.removeItem('token'); navigate('/login'); return; }
      const data = await res.json();
      setResults(data.results || []);
    } catch {
      setResults([]);
    } finally {
      setSearching(false);
    }
  }, [query, navigate]);

  const openPreview = async (result: MarketplaceResult) => {
    setSelectedResult(result);
    setLoadingPreview(true);
    setPreview(null);
    setInstallMsg('');
    try {
      const res = await fetch(`/api/v1/marketplace/preview/${result.owner}/${result.name}`, {
        headers: authHeader(),
      });
      const data = await res.json();
      setPreview(data);
    } catch {
      setPreview(null);
    } finally {
      setLoadingPreview(false);
    }
  };

  const install = async () => {
    if (!selectedResult) return;
    setInstalling(true);
    setInstallMsg('');
    try {
      const res = await fetch('/api/v1/marketplace/install', {
        method: 'POST',
        headers: { ...authHeader(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ owner: selectedResult.owner, repo: selectedResult.name }),
      });
      const data = await res.json();
      if (res.ok) {
        setInstallMsg(`✅ Template "${data.name}" v${data.version} installed successfully!`);
      } else {
        setInstallMsg(`❌ ${data.error || 'Installation failed.'}`);
      }
    } catch {
      setInstallMsg('❌ Network error during installation.');
    } finally {
      setInstalling(false);
    }
  };

  return (
    <div className="w-full max-w-5xl mx-auto p-6 animate-fade-in" style={{ margin: '0 auto' }}>
      {/* Header */}
      <div className="flex items-center gap-4 mb-8">
        <button onClick={() => navigate('/')} className="text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors">
          <ArrowLeft size={20} />
        </button>
        <div>
          <h1 className="text-3xl font-bold">Template Marketplace</h1>
          <p className="text-[var(--text-muted)] mt-1">Discover and install community service templates from GitHub</p>
        </div>
      </div>

      {/* Search Bar */}
      <Card glass className="mb-8">
        <div className="flex gap-3">
          <div style={{ position: 'relative', flex: 1 }}>
            <Search size={18} style={{ position: 'absolute', left: '14px', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }} />
            <input
              type="text"
              value={query}
              onChange={e => setQuery(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && search()}
              placeholder="Search templates (e.g. redis, postgres, nginx)..."
              style={{
                width: '100%',
                padding: '12px 14px 12px 44px',
                background: 'var(--surface)',
                border: '1px solid var(--border)',
                borderRadius: '10px',
                color: 'var(--text-primary)',
                fontSize: '15px',
              }}
            />
          </div>
          <Button onClick={search} disabled={searching || !query.trim()}>
            {searching ? 'Searching…' : 'Search'}
          </Button>
        </div>
      </Card>

      {/* Results */}
      {searching && (
        <div className="text-center text-[var(--text-muted)] py-12">Searching GitHub marketplace…</div>
      )}

      {!searching && searched && results.length === 0 && (
        <Card glass className="text-center py-12">
          <p className="text-[var(--text-muted)]">No templates found for "{query}".</p>
          <p className="text-[var(--text-muted)] text-sm mt-2">Try a different keyword or browse on GitHub.</p>
        </Card>
      )}

      {!searching && results.length > 0 && (
        <div style={{ display: 'grid', gap: '16px' }}>
          {results.map(r => (
            <Card key={`${r.owner}/${r.name}`} glass style={{ cursor: 'pointer' }} onClick={() => openPreview(r)}>
              <div className="flex justify-between items-start">
                <div style={{ flex: 1 }}>
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="font-bold text-lg">{r.name}</h3>
                    {r.is_official ? (
                      <span style={{
                        display: 'inline-flex', alignItems: 'center', gap: '4px',
                        padding: '2px 10px', borderRadius: '20px',
                        background: 'linear-gradient(135deg, #f59e0b, #d97706)',
                        color: '#fff', fontSize: '12px', fontWeight: 700,
                      }}>
                        <ShieldCheck size={12} /> Official
                      </span>
                    ) : (
                      <span style={{
                        display: 'inline-flex', alignItems: 'center', gap: '4px',
                        padding: '2px 10px', borderRadius: '20px',
                        background: 'rgba(99,102,241,0.2)',
                        color: 'var(--accent-primary)', fontSize: '12px', fontWeight: 600,
                      }}>
                        <Users size={12} /> Community
                      </span>
                    )}
                  </div>
                  <p className="text-[var(--text-muted)] text-sm mb-3">{r.description || 'No description available.'}</p>
                  <div className="flex items-center gap-4 text-sm text-[var(--text-secondary)]">
                    <span>by {r.owner}</span>
                    <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                      <Star size={14} className="text-yellow-400" /> {r.stars.toLocaleString()}
                    </span>
                  </div>
                </div>
                <div className="flex items-center gap-2 ml-4">
                  <a href={r.repo_url} target="_blank" rel="noreferrer" onClick={e => e.stopPropagation()}
                    className="text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors">
                    <ExternalLink size={16} />
                  </a>
                  <ChevronRight size={20} className="text-[var(--text-muted)]" />
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Preview Modal */}
      {selectedResult && (
        <div style={{
          position: 'fixed', inset: 0, zIndex: 50,
          background: 'rgba(0,0,0,0.7)', backdropFilter: 'blur(8px)',
          display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '24px',
        }} onClick={() => setSelectedResult(null)}>
          <div onClick={e => e.stopPropagation()} style={{
            background: 'var(--surface-elevated)', border: '1px solid var(--border)',
            borderRadius: '16px', width: '100%', maxWidth: '600px',
            maxHeight: '80vh', overflow: 'auto', padding: '28px',
          }}>
            <div className="flex justify-between items-start mb-6">
              <div>
                <div className="flex items-center gap-3 mb-1">
                  <h2 className="text-2xl font-bold">{selectedResult.name}</h2>
                  {selectedResult.is_official ? (
                    <span style={{ padding: '2px 10px', borderRadius: '20px', background: 'linear-gradient(135deg, #f59e0b, #d97706)', color: '#fff', fontSize: '12px', fontWeight: 700 }}>
                      ✅ Official
                    </span>
                  ) : (
                    <span style={{ padding: '2px 10px', borderRadius: '20px', background: 'rgba(239,68,68,0.15)', color: '#ef4444', fontSize: '12px', fontWeight: 600 }}>
                      ⚠️ Community — install at your own risk
                    </span>
                  )}
                </div>
                <p className="text-[var(--text-muted)] text-sm">by {selectedResult.owner}</p>
              </div>
              <button onClick={() => setSelectedResult(null)} className="text-[var(--text-muted)] hover:text-[var(--text-primary)]">
                <X size={20} />
              </button>
            </div>

            {loadingPreview && <div className="text-center text-[var(--text-muted)] py-8">Loading template details…</div>}

            {preview && (
              <>
                <div className="mb-4">
                  <p className="text-[var(--text-secondary)] mb-4">{preview.template.meta.description}</p>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', marginBottom: '16px' }}>
                    <div style={{ padding: '12px', background: 'var(--surface)', borderRadius: '10px' }}>
                      <div className="text-xs text-[var(--text-muted)] mb-1">Version</div>
                      <div className="font-medium">{preview.template.meta.version}</div>
                    </div>
                    <div style={{ padding: '12px', background: 'var(--surface)', borderRadius: '10px' }}>
                      <div className="text-xs text-[var(--text-muted)] mb-1">Image</div>
                      <div className="font-medium font-mono text-sm">{preview.template.runtime.image}</div>
                    </div>
                    <div style={{ padding: '12px', background: 'var(--surface)', borderRadius: '10px' }}>
                      <div className="text-xs text-[var(--text-muted)] mb-1">Supported Runtimes</div>
                      <div className="font-medium">{preview.template.runtime.supported?.join(', ')}</div>
                    </div>
                    <div style={{ padding: '12px', background: 'var(--surface)', borderRadius: '10px' }}>
                      <div className="text-xs text-[var(--text-muted)] mb-1">Tags</div>
                      <div className="font-medium">{preview.template.meta.tags?.join(', ') || '—'}</div>
                    </div>
                  </div>

                  {preview.template.config?.length > 0 && (
                    <div>
                      <h4 className="font-semibold mb-3">Configuration Variables</h4>
                      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                        {preview.template.config.map(c => (
                          <div key={c.name} style={{ padding: '10px 12px', background: 'var(--surface)', borderRadius: '8px', borderLeft: `3px solid ${c.required ? 'var(--error)' : 'var(--accent-primary)'}` }}>
                            <div className="flex justify-between items-center">
                              <code style={{ fontFamily: 'monospace', fontSize: '13px', color: 'var(--accent-primary)' }}>{c.name}</code>
                              <span style={{ fontSize: '11px', padding: '1px 6px', borderRadius: '4px', background: 'var(--surface-elevated)', color: 'var(--text-muted)' }}>
                                {c.type}{c.required ? ' · required' : ''}
                              </span>
                            </div>
                            {c.description && <p className="text-xs text-[var(--text-muted)] mt-1">{c.description}</p>}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>

                {installMsg && (
                  <div style={{ padding: '12px', borderRadius: '10px', background: installMsg.startsWith('✅') ? 'rgba(34,197,94,0.1)' : 'rgba(239,68,68,0.1)', marginBottom: '16px' }}>
                    {installMsg}
                  </div>
                )}

                <div className="flex gap-3 justify-end">
                  <Button variant="secondary" onClick={() => setSelectedResult(null)}>Cancel</Button>
                  <Button onClick={install} disabled={installing || installMsg.startsWith('✅')}>
                    <Download size={16} />
                    {installing ? 'Installing…' : 'Install Template'}
                  </Button>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
