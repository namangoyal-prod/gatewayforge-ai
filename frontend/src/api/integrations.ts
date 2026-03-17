import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Integration {
  id: string;
  name: string;
  partner_name: string;
  integration_type: string;
  payment_methods: string[];
  geographies: string[];
  expected_gmv: number;
  status: string;
  priority: string;
  created_by: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
  metadata: Record<string, any>;
}

export interface IntegrationStatus {
  integration: Integration;
  brd?: any;
  prd?: any;
  code?: any;
  test?: any;
  deployment?: any;
}

export interface PipelineMetric {
  id: string;
  integration_id: string;
  stage: string;
  started_at: string;
  completed_at?: string;
  duration_seconds?: number;
  status: string;
  ai_tokens_used?: number;
  ai_cost_usd?: number;
  errors?: any;
}

export const fetchIntegrations = async (): Promise<Integration[]> => {
  const response = await api.get('/integrations');
  return response.data.integrations;
};

export const fetchIntegration = async (id: string): Promise<Integration> => {
  const response = await api.get(`/integrations/${id}`);
  return response.data;
};

export const createIntegration = async (data: any): Promise<Integration> => {
  const response = await api.post('/integrations', data);
  return response.data;
};

export const updateIntegration = async (id: string, data: any): Promise<Integration> => {
  const response = await api.put(`/integrations/${id}`, data);
  return response.data;
};

export const deleteIntegration = async (id: string): Promise<void> => {
  await api.delete(`/integrations/${id}`);
};

export const fetchIntegrationStatus = async (id: string): Promise<IntegrationStatus> => {
  const response = await api.get(`/integrations/${id}/status`);
  return response.data;
};

export const fetchIntegrationTimeline = async (id: string): Promise<PipelineMetric[]> => {
  const response = await api.get(`/integrations/${id}/timeline`);
  return response.data.timeline;
};

export default api;
