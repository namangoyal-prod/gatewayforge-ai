import api from './integrations';

export interface BRDDocument {
  id: string;
  integration_id: string;
  file_name: string;
  file_path: string;
  file_type: string;
  uploaded_by: string;
  uploaded_at: string;
  validation_score?: number;
  validation_status: string;
  gap_analysis?: any;
  auto_fix_suggestions?: any;
  validated_at?: string;
  validated_by?: string;
}

export const uploadBRD = async (formData: FormData): Promise<BRDDocument> => {
  const response = await api.post('/brds', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};

export const fetchBRD = async (id: string): Promise<BRDDocument> => {
  const response = await api.get(`/brds/${id}`);
  return response.data;
};

export const validateBRD = async (id: string): Promise<void> => {
  await api.post(`/brds/${id}/validate`);
};

export const approveBRD = async (id: string, validatedBy: string): Promise<BRDDocument> => {
  const response = await api.post(`/brds/${id}/approve`, { validated_by: validatedBy });
  return response.data;
};

export const rejectBRD = async (id: string, validatedBy: string): Promise<BRDDocument> => {
  const response = await api.post(`/brds/${id}/reject`, { validated_by: validatedBy });
  return response.data;
};

export const fetchBRDGapAnalysis = async (id: string): Promise<any> => {
  const response = await api.get(`/brds/${id}/gap-analysis`);
  return response.data;
};
