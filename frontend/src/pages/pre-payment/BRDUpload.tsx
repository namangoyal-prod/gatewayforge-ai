import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { useDropzone } from 'react-dropzone';
import { uploadBRD } from '../api/brds';

const PAYMENT_METHODS = ['card', 'upi', 'netbanking', 'wallet'];
const GEOGRAPHIES = ['India', 'UAE', 'Singapore', 'Malaysia'];
const INTEGRATION_TYPES = [
  { value: 'gateway', label: 'Payment Gateway' },
  { value: 'aggregator', label: 'Aggregator' },
  { value: 'direct', label: 'Direct Integration' },
];

const BRDUpload = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    partner_name: '',
    integration_name: '',
    integration_type: 'gateway',
    payment_methods: [] as string[],
    geographies: [] as string[],
    expected_gmv: '',
    uploaded_by: 'current.user@razorpay.com',
  });
  const [file, setFile] = useState<File | null>(null);

  const { mutate: upload, isPending, error } = useMutation({
    mutationFn: uploadBRD,
    onSuccess: (data: any) => navigate(`/brd/${data.id}/validation`),
  });

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept: {
      'application/pdf': ['.pdf'],
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document': ['.docx'],
    },
    maxFiles: 1,
    onDrop: (files) => files.length > 0 && setFile(files[0]),
  });

  const toggle = (key: 'payment_methods' | 'geographies', val: string) => {
    const arr = formData[key];
    setFormData({ ...formData, [key]: arr.includes(val) ? arr.filter(v => v !== val) : [...arr, val] });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return alert('Please upload a BRD document');
    const fd = new FormData();
    fd.append('brd', file);
    Object.entries(formData).forEach(([k, v]) =>
      fd.append(k, Array.isArray(v) ? JSON.stringify(v) : v)
    );
    upload(fd);
  };

  const inputStyle: React.CSSProperties = {
    width: '100%', padding: '10px 12px', borderRadius: 8, fontSize: 14,
    border: '1px solid #d1d5db', outline: 'none', boxSizing: 'border-box',
  };

  const chipStyle = (active: boolean): React.CSSProperties => ({
    padding: '6px 14px', borderRadius: 20, fontSize: 13, cursor: 'pointer',
    border: `1px solid ${active ? '#6366f1' : '#d1d5db'}`,
    background: active ? '#ede9fe' : '#fff',
    color: active ? '#6366f1' : '#374151',
    fontWeight: active ? 600 : 400,
  });

  return (
    <div style={{ maxWidth: 720, margin: '0 auto' }}>
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 24, fontWeight: 700 }}>Upload BRD Document</h1>
        <p style={{ margin: '8px 0 0', color: '#6b7280', fontSize: 14 }}>
          Upload a Business Requirements Document for validation and PRD generation
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        <div style={{ background: '#fff', borderRadius: 12, padding: 28, boxShadow: '0 1px 3px rgba(0,0,0,0.1)', display: 'flex', flexDirection: 'column', gap: 24 }}>

          {/* Partner Info */}
          <section>
            <h2 style={{ margin: '0 0 16px', fontSize: 16, fontWeight: 600 }}>Partner Information</h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
              <div>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 500, marginBottom: 6 }}>Partner Name *</label>
                <input style={inputStyle} placeholder="e.g., HDFC Bank" value={formData.partner_name}
                  onChange={e => setFormData({ ...formData, partner_name: e.target.value })} required />
              </div>
              <div>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 500, marginBottom: 6 }}>Integration Name *</label>
                <input style={inputStyle} placeholder="e.g., HDFC UPI Gateway" value={formData.integration_name}
                  onChange={e => setFormData({ ...formData, integration_name: e.target.value })} required />
              </div>
              <div>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 500, marginBottom: 6 }}>Integration Type *</label>
                <select style={inputStyle} value={formData.integration_type}
                  onChange={e => setFormData({ ...formData, integration_type: e.target.value })}>
                  {INTEGRATION_TYPES.map(t => <option key={t.value} value={t.value}>{t.label}</option>)}
                </select>
              </div>
              <div>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 500, marginBottom: 8 }}>Payment Methods</label>
                <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                  {PAYMENT_METHODS.map(m => (
                    <button type="button" key={m} onClick={() => toggle('payment_methods', m)}
                      style={chipStyle(formData.payment_methods.includes(m))}>{m}</button>
                  ))}
                </div>
              </div>
              <div>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 500, marginBottom: 8 }}>Geographies</label>
                <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                  {GEOGRAPHIES.map(g => (
                    <button type="button" key={g} onClick={() => toggle('geographies', g)}
                      style={chipStyle(formData.geographies.includes(g))}>{g}</button>
                  ))}
                </div>
              </div>
              <div>
                <label style={{ display: 'block', fontSize: 13, fontWeight: 500, marginBottom: 6 }}>Expected GMV (₹ Cr)</label>
                <input style={inputStyle} type="number" placeholder="e.g., 500" value={formData.expected_gmv}
                  onChange={e => setFormData({ ...formData, expected_gmv: e.target.value })} />
              </div>
            </div>
          </section>

          {/* File Upload */}
          <section>
            <h2 style={{ margin: '0 0 16px', fontSize: 16, fontWeight: 600 }}>BRD Document</h2>
            <div {...getRootProps()} style={{
              padding: 40, textAlign: 'center', borderRadius: 10, cursor: 'pointer',
              border: `2px dashed ${isDragActive ? '#6366f1' : '#d1d5db'}`,
              background: isDragActive ? '#f5f3ff' : '#fafafa',
              transition: 'all 0.2s',
            }}>
              <input {...getInputProps()} />
              <div style={{ fontSize: 36, marginBottom: 12 }}>{file ? '📄' : '⬆️'}</div>
              {file ? (
                <>
                  <div style={{ fontWeight: 600, fontSize: 15 }}>{file.name}</div>
                  <div style={{ color: '#6b7280', fontSize: 13, marginTop: 4 }}>{(file.size / 1024 / 1024).toFixed(2)} MB</div>
                  <div style={{ color: '#10b981', fontSize: 13, marginTop: 6 }}>Click or drag to replace</div>
                </>
              ) : (
                <>
                  <div style={{ fontWeight: 600, fontSize: 15 }}>{isDragActive ? 'Drop here' : 'Drag & drop your BRD here'}</div>
                  <div style={{ color: '#6b7280', fontSize: 13, marginTop: 4 }}>or click to browse</div>
                  <div style={{ color: '#9ca3af', fontSize: 12, marginTop: 8 }}>Supports: PDF, DOCX (Max 50MB)</div>
                </>
              )}
            </div>
          </section>

          {/* Error */}
          {error && (
            <div style={{ background: '#fef2f2', border: '1px solid #fecaca', borderRadius: 8, padding: '12px 16px', color: '#dc2626', fontSize: 14 }}>
              {(error as Error).message}
            </div>
          )}

          {/* Actions */}
          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12 }}>
            <button type="button" onClick={() => navigate('/dashboard')} disabled={isPending}
              style={{ padding: '10px 20px', borderRadius: 8, border: '1px solid #d1d5db', background: '#fff', cursor: 'pointer', fontSize: 14 }}>
              Cancel
            </button>
            <button type="submit" disabled={!file || isPending}
              style={{ padding: '10px 20px', borderRadius: 8, border: 'none', background: !file || isPending ? '#c7d2fe' : '#6366f1', color: '#fff', cursor: !file || isPending ? 'not-allowed' : 'pointer', fontWeight: 600, fontSize: 14 }}>
              {isPending ? 'Uploading…' : '⬆️ Upload & Validate'}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default BRDUpload;
