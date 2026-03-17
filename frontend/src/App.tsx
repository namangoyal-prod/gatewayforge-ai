import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BladeProvider } from '@razorpay/blade/components';
import { bladeTheme } from '@razorpay/blade/tokens';

// Gateway pages
import Dashboard from './pages/gateway/Dashboard';
import IntegrationDetails from './pages/gateway/IntegrationDetails';

// Pre-payment processing pages
import BRDUpload from './pages/pre-payment/BRDUpload';
import BRDValidation from './pages/pre-payment/BRDValidation';
import EstimationEngine from './pages/pre-payment/EstimationEngine';

// Payment processing pages
import PRDReview from './pages/payment-processing/PRDReview';
import CodeReview from './pages/payment-processing/CodeReview';

// Post-payment processing pages
import TestExecution from './pages/post-payment/TestExecution';
import Deployment from './pages/post-payment/Deployment';
import Analytics from './pages/post-payment/Analytics';

// Layout
import AppLayout from './components/Layout/AppLayout';

// Create React Query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
});

function App() {
  return (
    <BladeProvider themeTokens={bladeTheme} colorScheme="light">
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<AppLayout />}>
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="dashboard" element={<Dashboard />} />
              <Route path="integrations/:id" element={<IntegrationDetails />} />
              <Route path="brd/upload" element={<BRDUpload />} />
              <Route path="brd/:id/validation" element={<BRDValidation />} />
              <Route path="prd/:id" element={<PRDReview />} />
              <Route path="code/:id" element={<CodeReview />} />
              <Route path="tests/:id" element={<TestExecution />} />
              <Route path="deployment/:id" element={<Deployment />} />
              <Route path="analytics" element={<Analytics />} />
              <Route path="estimation" element={<EstimationEngine />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </BladeProvider>
  );
}

export default App;
