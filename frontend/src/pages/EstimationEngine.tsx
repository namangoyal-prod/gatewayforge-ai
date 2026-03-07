import { useState } from 'react';
import api from '../api/integrations';

// ── Types ──────────────────────────────────────────────────────────────────────

interface MilestoneEstimate {
  name: string;
  category: string;
  traditional_days: number;
  ai_days: number;
  savings_days: number;
  savings_pct: number;
  ai_skills_used: string[];
  notes?: string;
  phase: number; // 1 | 2 | 0 (buffer)
}

interface EstimationResponse {
  bank_name: string;
  integration_type: string;
  codebase_reuse_label: string;
  milestones: MilestoneEstimate[];
  phase1_traditional_days: number;
  phase1_ai_days: number;
  phase2_traditional_days: number;
  phase2_ai_days: number;
  buffer_traditional_days: number;
  buffer_ai_days: number;
  total_traditional_days: number;
  total_ai_days: number;
  total_savings_days: number;
  total_savings_pct: number;
  references: ReferenceCalib[];
  assumptions: string[];
}

interface ReferenceCalib {
  name: string;
  traditional_days: number;
  ai_days: number;
  savings_pct: number;
  integration_type: string;
  notes: string;
}

// ── Constants ──────────────────────────────────────────────────────────────────

const CATEGORY_COLORS: Record<string, string> = {
  'Gateway Integration': '#6366f1',
  'Onboarding Changes':  '#f59e0b',
  'Merchant Dashboard':  '#06b6d4',
  'SoundBox & Devices':  '#8b5cf6',
  'QR Sync App':         '#ec4899',
  'Razorpay App':        '#10b981',
  'Reporting':           '#3b82f6',
  'NPCI Certification':  '#ef4444',
  'End-to-End Testing':  '#64748b',
  'Buffer':              '#9ca3af',
};

const PHASE_ACCENT: Record<number, string> = {
  1: '#6366f1',
  2: '#0ea5e9',
  0: '#9ca3af',
};

// ── Sub-components ─────────────────────────────────────────────────────────────

function Card({ children, style }: { children: React.ReactNode; style?: React.CSSProperties }) {
  return (
    <div style={{
      background: '#fff', borderRadius: 12, padding: 20,
      boxShadow: '0 1px 4px rgba(0,0,0,0.08)', ...style,
    }}>
      {children}
    </div>
  );
}

function StatCard({ label, value, sub, accent, dim }: {
  label: string; value: string; sub?: string; accent: string; dim?: boolean;
}) {
  return (
    <div style={{
      background: dim ? '#f9fafb' : '#fff',
      borderRadius: 12, padding: '16px 20px',
      boxShadow: dim ? 'none' : '0 1px 4px rgba(0,0,0,0.08)',
      borderTop: `4px solid ${accent}`,
      border: dim ? `1px solid #e5e7eb` : undefined,
      borderTopColor: accent, borderTopWidth: 4, borderTopStyle: 'solid',
    }}>
      <div style={{ fontSize: 11, color: '#6b7280', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: 6 }}>
        {label}
      </div>
      <div style={{ fontSize: 26, fontWeight: 700, color: dim ? '#6b7280' : '#111827' }}>{value}</div>
      {sub && <div style={{ fontSize: 12, color: '#9ca3af', marginTop: 4 }}>{sub}</div>}
    </div>
  );
}

function DayBar({ trad, ai, max }: { trad: number; ai: number; max: number }) {
  const w = (v: number) => Math.max(3, (v / max) * 140);
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
        <div style={{ width: w(trad), height: 7, borderRadius: 4, background: '#d1d5db', flexShrink: 0 }} />
        <span style={{ fontSize: 11, color: '#9ca3af', minWidth: 28 }}>{trad}d</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
        <div style={{ width: w(ai), height: 7, borderRadius: 4, background: '#6366f1', flexShrink: 0 }} />
        <span style={{ fontSize: 11, color: '#6366f1', minWidth: 28 }}>{ai}d</span>
      </div>
    </div>
  );
}

function PhaseTag({ phase }: { phase: number }) {
  if (phase === 0) return null;
  return (
    <span style={{
      fontSize: 10, fontWeight: 700, padding: '1px 6px', borderRadius: 4,
      background: PHASE_ACCENT[phase] + '18', color: PHASE_ACCENT[phase],
      marginRight: 6, verticalAlign: 'middle',
    }}>
      P{phase}
    </span>
  );
}

function SkillPill({ skill }: { skill: string }) {
  return (
    <span style={{
      fontSize: 10, padding: '2px 6px', borderRadius: 4,
      background: '#f3f4f6', color: '#374151', fontFamily: 'monospace',
      marginRight: 4, marginBottom: 2, display: 'inline-block',
    }}>
      {skill}
    </span>
  );
}

function SectionHeader({ label, trad, ai, color, phase }: {
  label: string; trad: number; ai: number; color: string; phase?: number;
}) {
  return (
    <div style={{
      display: 'flex', alignItems: 'center', gap: 8,
      padding: '10px 18px',
      background: color + '0e',
      borderBottom: '1px solid #f3f4f6',
    }}>
      <div style={{ width: 10, height: 10, borderRadius: '50%', background: color }} />
      <span style={{ fontWeight: 600, fontSize: 13, flex: 1 }}>{label}</span>
      {phase !== undefined && <PhaseTag phase={phase} />}
      <span style={{ fontSize: 12, color: '#6b7280' }}>
        {trad}d →{' '}
        <strong style={{ color: '#6366f1' }}>{ai}d</strong>
      </span>
    </div>
  );
}

// ── Toggle button ──────────────────────────────────────────────────────────────

function MethodBtn({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
    <button onClick={onClick} style={{
      padding: '5px 14px', borderRadius: 20, fontSize: 13, cursor: 'pointer',
      border: `2px solid ${active ? '#6366f1' : '#e5e7eb'}`,
      background: active ? '#eef2ff' : '#fff',
      color: active ? '#4f46e5' : '#374151',
      fontWeight: active ? 600 : 400,
    }}>
      {label}
    </button>
  );
}

// ── Main component ─────────────────────────────────────────────────────────────

export default function EstimationEngine() {
  const [form, setForm] = useState({
    bank_name: '',
    integration_type: 'new_gateway',
    payment_methods: [] as string[],
    encryption_type: 'aes_cbc',
    onboarding_type: 'both',
    has_devices: false,
    has_mobile_app: false,
    has_npci_cert: false,
    // Gap 1
    codebase_reuse: 'none',
    // Gap 2
    buffer_pct: 0,
    // Gap 3
    num_phases: 1,
    phase2_payment_methods: [] as string[],
  });

  const [result, setResult] = useState<EstimationResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [expanded, setExpanded] = useState<string | null>(null);

  const toggleMethod = (m: string, key: 'payment_methods' | 'phase2_payment_methods') =>
    setForm(f => ({
      ...f,
      [key]: (f[key] as string[]).includes(m)
        ? (f[key] as string[]).filter(x => x !== m)
        : [...(f[key] as string[]), m],
    }));

  const handleEstimate = async () => {
    if (!form.bank_name.trim()) { setError('Bank name is required'); return; }
    if (form.payment_methods.length === 0) { setError('Select at least one payment method'); return; }
    setError('');
    setLoading(true);
    try {
      const res = await api.post<EstimationResponse>('/estimation/estimate', form);
      setResult(res.data);
    } catch (e: any) {
      setError(e?.response?.data?.error || 'Failed to generate estimate');
    } finally {
      setLoading(false);
    }
  };

  // Group milestones by phase → category
  const byPhase = (phase: number) => result?.milestones.filter(m => m.phase === phase) ?? [];
  const byCategory = (milestones: MilestoneEstimate[]) => {
    const map: Record<string, MilestoneEstimate[]> = {};
    for (const m of milestones) {
      map[m.category] = map[m.category] || [];
      map[m.category].push(m);
    }
    return map;
  };

  const maxDays = result
    ? Math.max(...result.milestones.map(m => m.traditional_days), 1)
    : 1;

  const p1 = byPhase(1);
  const p2 = byPhase(2);
  const buf = byPhase(0);

  const renderPhaseSection = (milestones: MilestoneEstimate[], phaseNum: number) => {
    const grouped = byCategory(milestones);
    return Object.entries(grouped).map(([cat, items]) => {
      const catTrad = r(items.reduce((s, m) => s + m.traditional_days, 0));
      const catAI = r(items.reduce((s, m) => s + m.ai_days, 0));
      return (
        <Card key={`${phaseNum}-${cat}`} style={{ marginBottom: 12, padding: 0, overflow: 'hidden' }}>
          <SectionHeader
            label={cat} trad={catTrad} ai={catAI}
            color={CATEGORY_COLORS[cat] || '#6b7280'}
            phase={phaseNum}
          />
          {items.map((m, idx) => {
            const key = `${phaseNum}-${cat}-${idx}`;
            const open = expanded === key;
            return (
              <div
                key={key}
                onClick={() => setExpanded(open ? null : key)}
                style={{
                  padding: '11px 18px',
                  borderBottom: idx < items.length - 1 ? '1px solid #f9fafb' : 'none',
                  cursor: 'pointer',
                  background: open ? '#f8f7ff' : '#fff',
                }}
              >
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 165px 68px', alignItems: 'center', gap: 10 }}>
                  <div style={{ fontSize: 13, fontWeight: 500 }}>{m.name}</div>
                  <DayBar trad={m.traditional_days} ai={m.ai_days} max={maxDays} />
                  <div style={{ textAlign: 'right' }}>
                    <span style={{
                      fontSize: 12, fontWeight: 700,
                      color: m.savings_pct > 0 ? '#10b981' : '#9ca3af',
                    }}>
                      {m.savings_pct > 0 ? `-${m.savings_pct}%` : 'unchanged'}
                    </span>
                  </div>
                </div>
                {open && (
                  <div style={{ marginTop: 10, paddingTop: 10, borderTop: '1px dashed #e5e7eb' }}>
                    {m.ai_skills_used?.length > 0 && (
                      <div style={{ marginBottom: 6 }}>
                        <span style={{ fontSize: 11, color: '#6b7280', marginRight: 6 }}>Claude skills:</span>
                        {m.ai_skills_used.map(s => <SkillPill key={s} skill={s} />)}
                      </div>
                    )}
                    {m.notes && (
                      <div style={{ fontSize: 12, color: '#6b7280', lineHeight: 1.5 }}>{m.notes}</div>
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </Card>
      );
    });
  };

  return (
    <div style={{ maxWidth: 1240, margin: '0 auto' }}>
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 22, fontWeight: 700 }}>Dev Time Estimator</h1>
        <p style={{ margin: '6px 0 0', fontSize: 13, color: '#6b7280' }}>
          AI-assisted vs traditional timeline — calibrated against IDFC Bank tech spec (Feb 2026) and
          Likith Reddy's actual Phase I/II task breakdown (IDFC Omni migration xlsx).
        </p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '350px 1fr', gap: 20, alignItems: 'start' }}>

        {/* ── Form ─────────────────────────────────────────────────────────── */}
        <Card style={{ position: 'sticky', top: 20 }}>
          <div style={{ fontWeight: 600, fontSize: 15, marginBottom: 16 }}>Integration Parameters</div>

          <Label>Bank / Partner name</Label>
          <input value={form.bank_name}
            onChange={e => setForm(f => ({ ...f, bank_name: e.target.value }))}
            placeholder="e.g. IDFC Bank, Axis Bank"
            style={inputStyle}
          />

          <Label>Integration type</Label>
          <select value={form.integration_type}
            onChange={e => setForm(f => ({ ...f, integration_type: e.target.value }))}
            style={inputStyle}>
            <option value="new_gateway">New gateway (0→1)</option>
            <option value="new_payment_method">New payment method on existing gateway</option>
            <option value="extension">Extension / config change</option>
          </select>

          <Label>Payment methods — Phase 1</Label>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, marginBottom: 4 }}>
            {['cards', 'upi', 'wallet', 'bnpl'].map(m => (
              <MethodBtn key={m} label={m.toUpperCase()}
                active={form.payment_methods.includes(m)}
                onClick={() => toggleMethod(m, 'payment_methods')} />
            ))}
          </div>

          <Label>Encryption / security</Label>
          <select value={form.encryption_type}
            onChange={e => setForm(f => ({ ...f, encryption_type: e.target.value }))}
            style={inputStyle}>
            <option value="none">None</option>
            <option value="aes_cbc">AES-256-CBC (e.g. IDFC UPI / Montran)</option>
            <option value="rsa">RSA / OAuth2 RS256</option>
            <option value="dukpt">DUKPT (card-present hardware)</option>
          </select>

          <Label>Merchant onboarding</Label>
          <select value={form.onboarding_type}
            onChange={e => setForm(f => ({ ...f, onboarding_type: e.target.value }))}
            style={inputStyle}>
            <option value="api_based">API-based</option>
            <option value="file_based">File-based (CSV / SFTP)</option>
            <option value="both">Both (WL + Fiserv + SFDC + NTB/ETB)</option>
          </select>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 9, marginTop: 12, marginBottom: 14 }}>
            {[
              { key: 'has_devices', label: 'SoundBox / QR / POS device integration' },
              { key: 'has_mobile_app', label: 'Razorpay mobile app changes' },
              { key: 'has_npci_cert', label: 'NPCI certification required' },
            ].map(({ key, label }) => (
              <label key={key} style={{ display: 'flex', alignItems: 'center', gap: 10, cursor: 'pointer', fontSize: 13 }}>
                <input type="checkbox" checked={(form as any)[key]}
                  onChange={e => setForm(f => ({ ...f, [key]: e.target.checked }))}
                  style={{ width: 15, height: 15, accentColor: '#6366f1' }} />
                {label}
              </label>
            ))}
          </div>

          {/* ── Gap 1: Codebase reuse ─────────────────────────────────── */}
          <div style={{ borderTop: '1px solid #f3f4f6', paddingTop: 14, marginBottom: 4 }}>
            <div style={{ fontSize: 11, fontWeight: 700, color: '#6366f1', textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: 10 }}>
              Gap 1 — Codebase reuse
            </div>
            <Label>Existing infrastructure</Label>
            <select value={form.codebase_reuse}
              onChange={e => setForm(f => ({ ...f, codebase_reuse: e.target.value }))}
              style={inputStyle}>
              <option value="none">Greenfield — no existing infra</option>
              <option value="partial">Partial — some patterns exist (e.g. Montran already live)</option>
              <option value="migration">Migration — significant existing infra (e.g. IDFC Omni)</option>
            </select>
            {form.codebase_reuse === 'migration' && (
              <div style={{ fontSize: 11, color: '#6b7280', marginTop: 4, lineHeight: 1.5 }}>
                Traditional baseline deflated to 40% of greenfield.
                Calibrated: Likith's actual 97d vs tool greenfield 240d.
              </div>
            )}
          </div>

          {/* ── Gap 2: Buffer ─────────────────────────────────────────── */}
          <div style={{ borderTop: '1px solid #f3f4f6', paddingTop: 14, marginBottom: 4 }}>
            <div style={{ fontSize: 11, fontWeight: 700, color: '#f59e0b', textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: 10 }}>
              Gap 2 — Buffer modeling
            </div>
            <Label>Adhoc / oncall / leaves buffer (%)</Label>
            <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <input type="range" min={0} max={30} step={5}
                value={form.buffer_pct}
                onChange={e => setForm(f => ({ ...f, buffer_pct: Number(e.target.value) }))}
                style={{ flex: 1, accentColor: '#f59e0b' }} />
              <span style={{
                minWidth: 36, textAlign: 'center', fontWeight: 700, fontSize: 14,
                color: form.buffer_pct > 0 ? '#f59e0b' : '#9ca3af',
              }}>
                {form.buffer_pct}%
              </span>
            </div>
            {form.buffer_pct > 0 && (
              <div style={{ fontSize: 11, color: '#6b7280', marginTop: 4 }}>
                AI reduces buffer by 40% — fewer rework loops, SLITs catch bugs earlier.
              </div>
            )}
          </div>

          {/* ── Gap 3: Phase-based ────────────────────────────────────── */}
          <div style={{ borderTop: '1px solid #f3f4f6', paddingTop: 14, marginBottom: 16 }}>
            <div style={{ fontSize: 11, fontWeight: 700, color: '#0ea5e9', textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: 10 }}>
              Gap 3 — Phase-based estimation
            </div>
            <Label>Number of phases</Label>
            <div style={{ display: 'flex', gap: 8, marginBottom: 10 }}>
              {[1, 2].map(n => (
                <button key={n} onClick={() => setForm(f => ({ ...f, num_phases: n }))}
                  style={{
                    flex: 1, padding: '8px 0', borderRadius: 8, cursor: 'pointer', fontSize: 13, fontWeight: 600,
                    border: `2px solid ${form.num_phases === n ? '#0ea5e9' : '#e5e7eb'}`,
                    background: form.num_phases === n ? '#e0f2fe' : '#fff',
                    color: form.num_phases === n ? '#0284c7' : '#6b7280',
                  }}>
                  Phase {n}{n === 2 ? ' (55% reuse)' : ''}
                </button>
              ))}
            </div>

            {form.num_phases === 2 && (
              <>
                <Label>Phase 2 payment methods (leave empty to inherit)</Label>
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                  {['cards', 'upi', 'wallet', 'bnpl'].map(m => (
                    <MethodBtn key={m} label={m.toUpperCase()}
                      active={form.phase2_payment_methods.includes(m)}
                      onClick={() => toggleMethod(m, 'phase2_payment_methods')} />
                  ))}
                </div>
                <div style={{ fontSize: 11, color: '#6b7280', marginTop: 6 }}>
                  Phase 2 reuses all Phase 1 patterns at 55% effort. No device/dashboard/mobile re-work.
                </div>
              </>
            )}
          </div>

          {error && <div style={{ color: '#ef4444', fontSize: 13, marginBottom: 10 }}>{error}</div>}

          <button onClick={handleEstimate} disabled={loading} style={{
            width: '100%', padding: '11px 0',
            background: loading ? '#a5b4fc' : '#6366f1',
            color: '#fff', border: 'none', borderRadius: 8,
            fontWeight: 600, fontSize: 14, cursor: loading ? 'not-allowed' : 'pointer',
          }}>
            {loading ? 'Calculating…' : 'Generate Estimate'}
          </button>
        </Card>

        {/* ── Results ─────────────────────────────────────────────────────── */}
        <div>
          {!result ? (
            <Card style={{ textAlign: 'center', padding: 60, color: '#9ca3af' }}>
              <div style={{ fontSize: 40, marginBottom: 12 }}>⏱</div>
              <div style={{ fontSize: 15 }}>Fill parameters and click Generate Estimate</div>
              <div style={{ fontSize: 13, marginTop: 6 }}>Now models codebase reuse · buffer · multi-phase</div>
            </Card>
          ) : (
            <>
              {/* Reuse label banner */}
              <div style={{
                background: '#f0fdf4', border: '1px solid #bbf7d0', borderRadius: 10,
                padding: '10px 16px', fontSize: 13, color: '#166534', marginBottom: 16,
              }}>
                <strong>Calibration mode:</strong> {result.codebase_reuse_label}
              </div>

              {/* ── Stat cards ──────────────────────────────────────────── */}
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(5, 1fr)', gap: 12, marginBottom: 18 }}>
                <StatCard label="Traditional total" value={`${result.total_traditional_days}d`}
                  sub={`≈ ${Math.ceil(result.total_traditional_days / 22)} months`} accent="#d1d5db" />
                <StatCard label="With Claude skills" value={`${result.total_ai_days}d`}
                  sub={`≈ ${Math.ceil(result.total_ai_days / 22)} months`} accent="#6366f1" />
                <StatCard label="Days saved" value={`${result.total_savings_days}d`}
                  sub={`${result.total_savings_pct}% faster`} accent="#10b981" />
                <StatCard label="Phase 1" dim
                  value={`${result.phase1_ai_days}d`}
                  sub={`vs ${result.phase1_traditional_days}d trad`} accent={PHASE_ACCENT[1]} />
                {result.phase2_traditional_days > 0
                  ? <StatCard label="Phase 2 (55% reuse)" dim
                      value={`${result.phase2_ai_days}d`}
                      sub={`vs ${result.phase2_traditional_days}d trad`} accent={PHASE_ACCENT[2]} />
                  : result.buffer_traditional_days > 0
                    ? <StatCard label={`Buffer (${form.buffer_pct}%)`} dim
                        value={`${result.buffer_ai_days}d`}
                        sub={`vs ${result.buffer_traditional_days}d trad`} accent="#f59e0b" />
                    : <StatCard label="Bank" dim value={result.bank_name || '—'}
                        sub={result.integration_type.replace(/_/g, ' ')} accent="#f59e0b" />
                }
              </div>

              {/* ── Phase summary bar ────────────────────────────────────── */}
              {(result.phase2_traditional_days > 0 || result.buffer_traditional_days > 0) && (
                <Card style={{ marginBottom: 14 }}>
                  <div style={{ fontWeight: 600, fontSize: 13, marginBottom: 12 }}>Timeline breakdown</div>
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: 10 }}>
                    {[
                      { label: 'Phase 1 work', trad: result.phase1_traditional_days, ai: result.phase1_ai_days, accent: PHASE_ACCENT[1] },
                      ...(result.phase2_traditional_days > 0 ? [{ label: 'Phase 2 work', trad: result.phase2_traditional_days, ai: result.phase2_ai_days, accent: PHASE_ACCENT[2] }] : []),
                      ...(result.buffer_traditional_days > 0 ? [{ label: `Buffer (${form.buffer_pct}%)`, trad: result.buffer_traditional_days, ai: result.buffer_ai_days, accent: '#f59e0b' }] : []),
                      { label: 'Total', trad: result.total_traditional_days, ai: result.total_ai_days, accent: '#10b981' },
                    ].map(({ label, trad, ai, accent }) => (
                      <div key={label} style={{ textAlign: 'center', padding: '10px 0' }}>
                        <div style={{ fontSize: 11, color: '#6b7280', marginBottom: 6 }}>{label}</div>
                        <div style={{ fontSize: 11, color: '#9ca3af', marginBottom: 2 }}>{trad}d trad</div>
                        <div style={{ fontSize: 20, fontWeight: 700, color: accent }}>{ai}d</div>
                        <div style={{ fontSize: 11, color: '#10b981', marginTop: 2 }}>
                          -{Math.round(((trad - ai) / trad) * 100)}%
                        </div>
                      </div>
                    ))}
                  </div>
                </Card>
              )}

              {/* Legend */}
              <div style={{ display: 'flex', gap: 20, marginBottom: 14, fontSize: 12, color: '#6b7280', flexWrap: 'wrap' }}>
                <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                  <span style={{ width: 20, height: 7, background: '#d1d5db', borderRadius: 4, display: 'inline-block' }} /> Traditional
                </span>
                <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                  <span style={{ width: 20, height: 7, background: '#6366f1', borderRadius: 4, display: 'inline-block' }} /> AI-assisted
                </span>
                {result.phase2_traditional_days > 0 && <>
                  <span style={{ display: 'flex', alignItems: 'center', gap: 5 }}>
                    <span style={{ width: 14, height: 14, background: '#eef2ff', border: '1.5px solid #6366f1', borderRadius: 3, display: 'inline-block' }} /> Phase 1
                  </span>
                  <span style={{ display: 'flex', alignItems: 'center', gap: 5 }}>
                    <span style={{ width: 14, height: 14, background: '#e0f2fe', border: '1.5px solid #0ea5e9', borderRadius: 3, display: 'inline-block' }} /> Phase 2
                  </span>
                </>}
                <span style={{ marginLeft: 'auto', fontSize: 11 }}>Click a row to expand skills + notes</span>
              </div>

              {/* ── Phase 1 milestones ─────────────────────────────────── */}
              {result.phase2_traditional_days > 0 && (
                <div style={{
                  fontSize: 12, fontWeight: 700, color: PHASE_ACCENT[1], textTransform: 'uppercase',
                  letterSpacing: '0.05em', marginBottom: 8, paddingLeft: 4,
                }}>
                  Phase 1 — {result.phase1_traditional_days}d → {result.phase1_ai_days}d with AI
                </div>
              )}
              {renderPhaseSection(p1, 1)}

              {/* ── Phase 2 milestones ─────────────────────────────────── */}
              {p2.length > 0 && (
                <>
                  <div style={{
                    fontSize: 12, fontWeight: 700, color: PHASE_ACCENT[2], textTransform: 'uppercase',
                    letterSpacing: '0.05em', margin: '16px 0 8px', paddingLeft: 4,
                  }}>
                    Phase 2 (55% pattern reuse) — {result.phase2_traditional_days}d → {result.phase2_ai_days}d with AI
                  </div>
                  {renderPhaseSection(p2, 2)}
                </>
              )}

              {/* ── Buffer ────────────────────────────────────────────────── */}
              {buf.length > 0 && (
                <>
                  <div style={{
                    fontSize: 12, fontWeight: 700, color: '#9ca3af', textTransform: 'uppercase',
                    letterSpacing: '0.05em', margin: '16px 0 8px', paddingLeft: 4,
                  }}>
                    Buffer — {result.buffer_traditional_days}d → {result.buffer_ai_days}d with AI
                  </div>
                  {buf.map((m, i) => {
                    const key = `buf-${i}`;
                    const open = expanded === key;
                    return (
                      <Card key={key} style={{ marginBottom: 12, padding: 0, overflow: 'hidden' }}>
                        <div style={{ padding: '11px 18px', cursor: 'pointer', background: open ? '#fafafa' : '#fff' }}
                          onClick={() => setExpanded(open ? null : key)}>
                          <div style={{ display: 'grid', gridTemplateColumns: '1fr 165px 68px', alignItems: 'center', gap: 10 }}>
                            <div style={{ fontSize: 13, fontWeight: 500, color: '#6b7280' }}>{m.name}</div>
                            <DayBar trad={m.traditional_days} ai={m.ai_days} max={maxDays} />
                            <div style={{ textAlign: 'right' }}>
                              <span style={{ fontSize: 12, fontWeight: 700, color: '#10b981' }}>-{m.savings_pct}%</span>
                            </div>
                          </div>
                          {open && m.notes && (
                            <div style={{ marginTop: 10, paddingTop: 10, borderTop: '1px dashed #e5e7eb', fontSize: 12, color: '#6b7280', lineHeight: 1.5 }}>
                              {m.notes}
                            </div>
                          )}
                        </div>
                      </Card>
                    );
                  })}
                </>
              )}

              {/* ── Reference calibrations ──────────────────────────────── */}
              <Card style={{ marginBottom: 14, marginTop: 8 }}>
                <div style={{ fontWeight: 600, fontSize: 14, marginBottom: 14 }}>Reference calibration data</div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
                  {result.references.map(ref => (
                    <div key={ref.name} style={{
                      padding: '12px 14px', borderRadius: 8, background: '#f9fafb',
                      borderLeft: `3px solid ${ref.integration_type === 'migration' ? '#10b981' : '#6366f1'}`,
                    }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <span style={{ fontWeight: 600, fontSize: 13 }}>{ref.name}</span>
                        <span style={{
                          fontSize: 12, fontWeight: 700, color: '#10b981',
                          background: '#d1fae5', padding: '2px 8px', borderRadius: 10,
                        }}>-{ref.savings_pct}%</span>
                      </div>
                      <div style={{ fontSize: 12, color: '#6b7280', marginTop: 4 }}>
                        {ref.traditional_days}d traditional → {ref.ai_days}d with AI · {ref.notes}
                      </div>
                    </div>
                  ))}
                </div>
              </Card>

              {/* ── Assumptions ─────────────────────────────────────────── */}
              <Card>
                <div style={{ fontWeight: 600, fontSize: 14, marginBottom: 12 }}>Assumptions</div>
                <ul style={{ margin: 0, paddingLeft: 18, display: 'flex', flexDirection: 'column', gap: 5 }}>
                  {result.assumptions.map((a, i) => (
                    <li key={i} style={{ fontSize: 12, color: '#6b7280', lineHeight: 1.5 }}>{a}</li>
                  ))}
                </ul>
              </Card>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

// ── Helpers ────────────────────────────────────────────────────────────────────

function Label({ children }: { children: React.ReactNode }) {
  return (
    <div style={{ fontSize: 12, fontWeight: 600, color: '#374151', marginBottom: 5, marginTop: 12 }}>
      {children}
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  width: '100%', padding: '8px 10px', borderRadius: 8,
  border: '1px solid #e5e7eb', fontSize: 13, outline: 'none',
  background: '#fafafa', boxSizing: 'border-box',
};

function r(n: number) { return Math.round(n * 10) / 10; }
