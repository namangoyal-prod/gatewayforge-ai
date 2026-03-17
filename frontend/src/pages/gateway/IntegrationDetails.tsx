import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchIntegrationStatus, fetchIntegrationTimeline } from '../api/integrations';
import { formatDistanceToNow } from 'date-fns';

const STAGES = ['BRD Validation', 'PRD Generation', 'Code Generation', 'Testing', 'Deployment'];

const STAGE_MAP: Record<string, number> = {
  brd_uploaded: 1, brd_validation: 1,
  prd_generation: 2, prd_approved: 2,
  code_generation: 3, code_approved: 3,
  testing: 4, tests_passed: 4,
  deploying: 5, deployed: 5,
};

const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  pending:  { bg: '#f3f4f6', text: '#6b7280' },
  running:  { bg: '#dbeafe', text: '#1d4ed8' },
  approved: { bg: '#d1fae5', text: '#065f46' },
  success:  { bg: '#d1fae5', text: '#065f46' },
  passed:   { bg: '#d1fae5', text: '#065f46' },
  failed:   { bg: '#fee2e2', text: '#b91c1c' },
  rejected: { bg: '#fee2e2', text: '#b91c1c' },
  deployed: { bg: '#d1fae5', text: '#065f46' },
};

const Badge = ({ status }: { status: string }) => {
  const c = STATUS_COLORS[status] || STATUS_COLORS.pending;
  return (
    <span style={{ background: c.bg, color: c.text, padding: '2px 10px', borderRadius: 12, fontSize: 12, fontWeight: 600 }}>
      {status}
    </span>
  );
};

const TABS = ['Overview', 'BRD', 'PRD', 'Code', 'Tests', 'Deployment', 'Timeline'];

const IntegrationDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState(0);

  const { data: status, isLoading } = useQuery({
    queryKey: ['integration-status', id],
    queryFn: () => fetchIntegrationStatus(id!),
    refetchInterval: 5000,
  });

  const { data: timeline } = useQuery({
    queryKey: ['integration-timeline', id],
    queryFn: () => fetchIntegrationTimeline(id!),
  });

  if (isLoading) return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 300, color: '#6b7280' }}>
      Loading integration details…
    </div>
  );

  if (!status) return <div style={{ padding: 24 }}>Integration not found.</div>;

  const integration = status.integration;
  const stageNum = STAGE_MAP[integration.status] || 0;
  const progress = Math.round((stageNum / 5) * 100);

  return (
    <div>
      {/* Back */}
      <button onClick={() => navigate('/dashboard')} style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#6366f1', fontSize: 14, marginBottom: 16, padding: 0 }}>
        ← Back to Dashboard
      </button>

      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 24 }}>
        <div>
          <h1 style={{ margin: 0, fontSize: 24, fontWeight: 700 }}>{integration.partner_name}</h1>
          <div style={{ color: '#6b7280', marginTop: 4, fontSize: 15 }}>{integration.name}</div>
          <div style={{ display: 'flex', gap: 8, marginTop: 10 }}>
            <Badge status={integration.integration_type} />
            <Badge status={integration.priority} />
            <Badge status={integration.status} />
          </div>
        </div>
        <div style={{ display: 'flex', gap: 10 }}>
          <button style={{ padding: '8px 16px', borderRadius: 8, border: '1px solid #d1d5db', background: '#fff', cursor: 'pointer', fontSize: 14 }}>View Logs</button>
          <button style={{ padding: '8px 16px', borderRadius: 8, border: 'none', background: '#6366f1', color: '#fff', cursor: 'pointer', fontWeight: 600, fontSize: 14 }}>Take Action</button>
        </div>
      </div>

      {/* Progress */}
      <div style={{ background: '#fff', borderRadius: 12, padding: 20, marginBottom: 24, boxShadow: '0 1px 3px rgba(0,0,0,0.08)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 10 }}>
          <span style={{ fontWeight: 600, fontSize: 14 }}>Pipeline Progress</span>
          <span style={{ color: '#6b7280', fontSize: 13 }}>{stageNum} of 5 stages completed</span>
        </div>
        <div style={{ background: '#f3f4f6', borderRadius: 99, height: 8, overflow: 'hidden' }}>
          <div style={{ width: `${progress}%`, height: '100%', background: progress === 100 ? '#10b981' : '#6366f1', borderRadius: 99, transition: 'width 0.4s' }} />
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 8 }}>
          {STAGES.map((s, i) => (
            <span key={s} style={{ fontSize: 11, color: i < stageNum ? '#6366f1' : '#9ca3af', fontWeight: i < stageNum ? 600 : 400 }}>{s}</span>
          ))}
        </div>
      </div>

      {/* Tabs */}
      <div style={{ borderBottom: '1px solid #e5e7eb', marginBottom: 24, display: 'flex', gap: 0 }}>
        {TABS.map((tab, i) => (
          <button key={tab} onClick={() => setActiveTab(i)} style={{
            padding: '10px 16px', border: 'none', background: 'none', cursor: 'pointer', fontSize: 14,
            borderBottom: activeTab === i ? '2px solid #6366f1' : '2px solid transparent',
            color: activeTab === i ? '#6366f1' : '#6b7280', fontWeight: activeTab === i ? 600 : 400,
          }}>{tab}</button>
        ))}
      </div>

      {/* Tab Content */}
      <div style={{ background: '#fff', borderRadius: 12, padding: 24, boxShadow: '0 1px 3px rgba(0,0,0,0.08)' }}>

        {/* Overview */}
        {activeTab === 0 && (
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 24 }}>
            <div>
              <h3 style={{ margin: '0 0 16px', fontSize: 15 }}>Integration Details</h3>
              {[
                ['Partner', integration.partner_name],
                ['Type', integration.integration_type],
                ['Payment Methods', integration.payment_methods?.join(', ') || 'N/A'],
                ['Geographies', integration.geographies?.join(', ') || 'N/A'],
                ['Expected GMV', `₹${integration.expected_gmv} Cr`],
                ['Created By', integration.created_by],
                ['Created', formatDistanceToNow(new Date(integration.created_at), { addSuffix: true })],
              ].map(([label, value]) => (
                <div key={label} style={{ marginBottom: 12 }}>
                  <div style={{ fontSize: 12, color: '#9ca3af' }}>{label}</div>
                  <div style={{ fontSize: 14, fontWeight: 500 }}>{value}</div>
                </div>
              ))}
            </div>
            <div>
              <h3 style={{ margin: '0 0 16px', fontSize: 15 }}>Stage Status</h3>
              {[
                ['BRD Validation', status.brd?.validation_status || 'pending'],
                ['PRD Generation', status.prd?.status || 'pending'],
                ['Code Generation', status.code?.code_review_status || 'pending'],
                ['Test Execution', status.test?.passed ? 'passed' : 'pending'],
                ['Deployment', status.deployment?.status || 'pending'],
              ].map(([label, s]) => (
                <div key={label} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                  <span style={{ fontSize: 14 }}>{label}</span>
                  <Badge status={s} />
                </div>
              ))}
            </div>
          </div>
        )}

        {/* BRD */}
        {activeTab === 1 && (
          status.brd ? (
            <div>
              <h3 style={{ margin: '0 0 16px', fontSize: 15 }}>BRD Validation Report — {status.brd.file_name}</h3>
              <div style={{ marginBottom: 16 }}>
                <div style={{ fontSize: 12, color: '#9ca3af' }}>Validation Score</div>
                <div style={{ fontSize: 40, fontWeight: 700, color: '#6366f1' }}>{status.brd.validation_score ?? '—'}<span style={{ fontSize: 16, color: '#9ca3af' }}>/100</span></div>
              </div>
              <div style={{ marginBottom: 16 }}>
                <div style={{ fontSize: 12, color: '#9ca3af' }}>Status</div>
                <Badge status={status.brd.validation_status} />
              </div>
              <button onClick={() => navigate(`/brd/${status.brd.id}/validation`)}
                style={{ padding: '10px 20px', borderRadius: 8, border: 'none', background: '#6366f1', color: '#fff', cursor: 'pointer', fontWeight: 600, fontSize: 14 }}>
                View Full Report
              </button>
            </div>
          ) : <div style={{ color: '#6b7280' }}>No BRD uploaded yet.</div>
        )}

        {/* PRD */}
        {activeTab === 2 && (
          status.prd ? (
            <div>
              <h3 style={{ margin: '0 0 12px', fontSize: 15 }}>Generated PRD</h3>
              <Badge status={status.prd.status} />
              <pre style={{ marginTop: 16, background: '#f9fafb', padding: 16, borderRadius: 8, fontSize: 13, whiteSpace: 'pre-wrap', maxHeight: 400, overflow: 'auto' }}>
                {status.prd.content}
              </pre>
            </div>
          ) : <div style={{ color: '#6b7280' }}>PRD not generated yet.</div>
        )}

        {/* Code */}
        {activeTab === 3 && (
          status.code ? (
            <div>
              <h3 style={{ margin: '0 0 12px', fontSize: 15 }}>Code Generation</h3>
              <div style={{ marginBottom: 8 }}><Badge status={status.code.code_review_status} /></div>
              <div style={{ fontSize: 13, color: '#6b7280' }}>Reference: {status.code.reference_integration || 'N/A'}</div>
            </div>
          ) : <div style={{ color: '#6b7280' }}>Code not generated yet.</div>
        )}

        {/* Tests */}
        {activeTab === 4 && (
          status.test ? (
            <div>
              <h3 style={{ margin: '0 0 16px', fontSize: 15 }}>Test Suite</h3>
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
                {[
                  ['Unit Tests', status.test.unit_tests_count],
                  ['Integration Tests', status.test.integration_tests_count],
                  ['E2E Tests', status.test.e2e_tests_count],
                  ['Coverage', `${status.test.total_coverage_percent ?? 0}%`],
                  ['Status', status.test.passed ? 'Passed ✅' : 'Pending'],
                ].map(([label, value]) => (
                  <div key={label} style={{ background: '#f9fafb', borderRadius: 8, padding: 16 }}>
                    <div style={{ fontSize: 12, color: '#9ca3af' }}>{label}</div>
                    <div style={{ fontSize: 20, fontWeight: 700 }}>{value}</div>
                  </div>
                ))}
              </div>
            </div>
          ) : <div style={{ color: '#6b7280' }}>Tests not generated yet.</div>
        )}

        {/* Deployment */}
        {activeTab === 5 && (
          status.deployment ? (
            <div>
              <h3 style={{ margin: '0 0 12px', fontSize: 15 }}>Deployment</h3>
              <Badge status={status.deployment.status} />
              <div style={{ marginTop: 12, fontSize: 13, color: '#6b7280' }}>Environment: {status.deployment.environment}</div>
            </div>
          ) : <div style={{ color: '#6b7280' }}>Not deployed yet.</div>
        )}

        {/* Timeline */}
        {activeTab === 6 && (
          <div>
            <h3 style={{ margin: '0 0 16px', fontSize: 15 }}>Pipeline Timeline</h3>
            {timeline && timeline.length > 0 ? (
              <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
                {timeline.map((m: any) => (
                  <div key={m.id} style={{ display: 'flex', gap: 14, alignItems: 'flex-start' }}>
                    <div style={{ fontSize: 20, marginTop: 2 }}>
                      {m.status === 'success' ? '✅' : m.status === 'failed' ? '❌' : '⏳'}
                    </div>
                    <div>
                      <div style={{ fontWeight: 600, fontSize: 14 }}>{m.stage}</div>
                      <div style={{ fontSize: 12, color: '#9ca3af', marginTop: 2 }}>
                        {m.status === 'success' ? 'Completed' : 'In Progress'}
                        {m.duration_seconds && ` · ${m.duration_seconds}s`}
                        {' · '}{formatDistanceToNow(new Date(m.started_at), { addSuffix: true })}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div style={{ color: '#6b7280' }}>No timeline events yet.</div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default IntegrationDetails;
