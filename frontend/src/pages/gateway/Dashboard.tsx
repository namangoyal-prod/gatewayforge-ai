import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchIntegrations } from '../api/integrations';

const PIPELINE_STAGES = [
  { id: 'brd_uploaded', label: 'BRD Review', color: '#3b82f6' },
  { id: 'prd_generation', label: 'PRD Generation', color: '#f59e0b' },
  { id: 'code_generation', label: 'Code Generation', color: '#8b5cf6' },
  { id: 'testing', label: 'Testing', color: '#06b6d4' },
  { id: 'deployed', label: 'Deployed', color: '#10b981' },
] as const;

const PRIORITY_COLORS: Record<string, string> = {
  critical: '#ef4444',
  high: '#f97316',
  medium: '#3b82f6',
  low: '#6b7280',
};

const Dashboard = () => {
  const navigate = useNavigate();

  const { data: integrations, isLoading, error } = useQuery({
    queryKey: ['integrations'],
    queryFn: fetchIntegrations,
    refetchInterval: 5000,
  });

  const grouped = PIPELINE_STAGES.reduce((acc, stage) => {
    acc[stage.id] = integrations?.filter(i => i.status === stage.id) || [];
    return acc;
  }, {} as Record<string, any[]>);

  const metrics = {
    total: integrations?.length || 0,
    inProgress: integrations?.filter(i => !['deployed', 'failed'].includes(i.status)).length || 0,
    completed: integrations?.filter(i => i.status === 'deployed').length || 0,
  };

  if (isLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 300, color: '#6b7280' }}>
        Loading integrations…
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: 16, color: '#ef4444' }}>
        Error: {(error as Error).message}
      </div>
    );
  }

  return (
    <div>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 24, fontWeight: 700 }}>Integration Pipeline</h1>
        <button
          onClick={() => navigate('/brd/upload')}
          style={{
            background: '#6366f1', color: '#fff', border: 'none', borderRadius: 8,
            padding: '10px 20px', cursor: 'pointer', fontWeight: 600, fontSize: 14,
          }}
        >
          + New Integration
        </button>
      </div>

      {/* Metrics */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16, marginBottom: 28 }}>
        {[
          { label: 'Total Integrations', value: metrics.total },
          { label: 'In Progress', value: metrics.inProgress },
          { label: 'Completed', value: metrics.completed },
        ].map(({ label, value }) => (
          <div key={label} style={{
            background: '#fff', borderRadius: 12, padding: '20px 24px',
            boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
          }}>
            <div style={{ fontSize: 13, color: '#6b7280', marginBottom: 8 }}>{label}</div>
            <div style={{ fontSize: 32, fontWeight: 700, color: '#111827' }}>{value}</div>
          </div>
        ))}
      </div>

      {/* Kanban Board */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(5, 1fr)', gap: 16, overflowX: 'auto' }}>
        {PIPELINE_STAGES.map(stage => (
          <div key={stage.id}>
            {/* Stage Header */}
            <div style={{
              padding: '10px 12px', borderRadius: 8, marginBottom: 12,
              background: '#fff', boxShadow: '0 1px 3px rgba(0,0,0,0.08)',
              display: 'flex', alignItems: 'center', gap: 8,
            }}>
              <div style={{ width: 10, height: 10, borderRadius: '50%', background: stage.color }} />
              <span style={{ fontWeight: 600, fontSize: 13, flex: 1 }}>{stage.label}</span>
              <span style={{
                background: '#e5e7eb', borderRadius: 12, padding: '2px 8px', fontSize: 12,
              }}>
                {grouped[stage.id]?.length || 0}
              </span>
            </div>

            {/* Cards */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              {grouped[stage.id]?.map(integration => (
                <div
                  key={integration.id}
                  onClick={() => navigate(`/integrations/${integration.id}`)}
                  style={{
                    background: '#fff', borderRadius: 10, padding: 14, cursor: 'pointer',
                    boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
                    transition: 'box-shadow 0.15s',
                    border: '1px solid #f3f4f6',
                  }}
                >
                  <div style={{ fontWeight: 600, fontSize: 14, marginBottom: 4 }}>
                    {integration.partner_name}
                  </div>
                  <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 8 }}>
                    {integration.name}
                  </div>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4, marginBottom: 8 }}>
                    {integration.payment_methods?.map((m: string) => (
                      <span key={m} style={{
                        background: '#f3f4f6', borderRadius: 4, padding: '2px 6px', fontSize: 11, color: '#374151',
                      }}>
                        {m}
                      </span>
                    ))}
                  </div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <span style={{ fontSize: 11, color: '#9ca3af' }}>
                      {new Date(integration.created_at).toLocaleDateString()}
                    </span>
                    <span style={{
                      fontSize: 11, fontWeight: 600, padding: '2px 8px', borderRadius: 4,
                      background: PRIORITY_COLORS[integration.priority] + '20',
                      color: PRIORITY_COLORS[integration.priority],
                    }}>
                      {integration.priority}
                    </span>
                  </div>
                </div>
              ))}

              {(!grouped[stage.id] || grouped[stage.id].length === 0) && (
                <div style={{
                  padding: 20, textAlign: 'center', color: '#d1d5db', fontSize: 13,
                  border: '2px dashed #e5e7eb', borderRadius: 10,
                }}>
                  No integrations
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Dashboard;
