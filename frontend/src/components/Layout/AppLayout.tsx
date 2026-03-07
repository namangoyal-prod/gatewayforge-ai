import { Outlet, useNavigate, useLocation } from 'react-router-dom';

const navItems = [
  { path: '/dashboard', label: 'Dashboard', emoji: '🏠' },
  { path: '/brd/upload', label: 'Upload BRD', emoji: '📄' },
  { path: '/analytics', label: 'Analytics', emoji: '📊' },
  { path: '/estimation', label: 'Dev Time Estimator', emoji: '⏱' },
];

const AppLayout = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const currentLabel = navItems.find(item => item.path === location.pathname)?.label || 'GatewayForge AI';

  return (
    <div style={{ display: 'flex', height: '100vh', fontFamily: 'system-ui, sans-serif' }}>
      {/* Sidebar */}
      <div style={{
        width: 240,
        background: '#1a1a2e',
        color: '#fff',
        display: 'flex',
        flexDirection: 'column',
        flexShrink: 0,
      }}>
        <div style={{ padding: '20px 16px', borderBottom: '1px solid rgba(255,255,255,0.1)' }}>
          <div style={{ fontWeight: 700, fontSize: 16 }}>GatewayForge AI</div>
          <div style={{ fontSize: 12, color: 'rgba(255,255,255,0.5)', marginTop: 4 }}>
            Autonomous Gateway Integration
          </div>
        </div>

        <nav style={{ flex: 1, padding: '12px 8px' }}>
          {navItems.map(item => (
            <button
              key={item.path}
              onClick={() => navigate(item.path)}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 10,
                width: '100%',
                padding: '10px 12px',
                marginBottom: 4,
                borderRadius: 8,
                border: 'none',
                cursor: 'pointer',
                background: location.pathname === item.path
                  ? 'rgba(99,102,241,0.3)'
                  : 'transparent',
                color: location.pathname === item.path ? '#a5b4fc' : 'rgba(255,255,255,0.7)',
                fontSize: 14,
                textAlign: 'left',
              }}
            >
              <span>{item.emoji}</span>
              <span>{item.label}</span>
            </button>
          ))}
        </nav>

        <div style={{ padding: '12px 16px', borderTop: '1px solid rgba(255,255,255,0.1)' }}>
          <div style={{ fontSize: 13, color: 'rgba(255,255,255,0.7)' }}>Product Manager</div>
          <div style={{ fontSize: 12, color: 'rgba(255,255,255,0.4)' }}>pm@razorpay.com</div>
        </div>
      </div>

      {/* Main Content */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* Top Bar */}
        <div style={{
          padding: '0 24px',
          height: 56,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          borderBottom: '1px solid #e5e7eb',
          background: '#fff',
          flexShrink: 0,
        }}>
          <span style={{ fontWeight: 600, fontSize: 16 }}>{currentLabel}</span>
          <span style={{ fontSize: 13, color: '#6b7280' }}>
            {new Date().toLocaleDateString('en-IN', { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' })}
          </span>
        </div>

        {/* Page Content */}
        <div style={{ flex: 1, overflow: 'auto', background: '#f9fafb', padding: 28 }}>
          <Outlet />
        </div>
      </div>
    </div>
  );
};

export default AppLayout;
