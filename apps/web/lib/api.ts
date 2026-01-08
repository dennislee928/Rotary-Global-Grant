import type {
  Report,
  ReportDetail,
  TriageDecision,
  Alert,
  TrainingEvent,
  QuizResult,
  TrainingStats,
  KPIMetrics,
  DashboardStats,
  PaginatedResponse,
  CreateReportRequest,
  TriageRequest,
  CreateAlertRequest,
} from './types';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || 'http://localhost:8080';

class ApiError extends Error {
  constructor(public code: string, message: string, public status: number) {
    super(message);
    this.name = 'ApiError';
  }
}

async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE}${endpoint}`;
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...options.headers as Record<string, string>,
  };

  // Add auth token if available
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ 
      code: 'UNKNOWN_ERROR', 
      message: 'An unknown error occurred' 
    }));
    throw new ApiError(error.code, error.message, response.status);
  }

  return response.json();
}

// Health check
export async function getHealth(): Promise<{
  status: string;
  ts: string;
  version: string;
  database: string;
  redis: string;
}> {
  return fetchApi('/healthz');
}

// Reports
export async function createReport(data: CreateReportRequest): Promise<Report> {
  return fetchApi('/v1/reports', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getReports(params?: {
  page?: number;
  pageSize?: number;
  status?: string;
  category?: string;
}): Promise<PaginatedResponse<Report>> {
  const searchParams = new URLSearchParams();
  if (params?.page) searchParams.set('page', params.page.toString());
  if (params?.pageSize) searchParams.set('pageSize', params.pageSize.toString());
  if (params?.status) searchParams.set('status', params.status);
  if (params?.category) searchParams.set('category', params.category);
  
  const query = searchParams.toString();
  return fetchApi(`/v1/reports${query ? `?${query}` : ''}`);
}

export async function getReport(id: string): Promise<ReportDetail> {
  return fetchApi(`/v1/reports/${id}`);
}

// Triage
export async function triageReport(
  reportId: string,
  data: TriageRequest
): Promise<TriageDecision> {
  return fetchApi(`/v1/reports/${reportId}/triage`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getTriageDecisions(params?: {
  page?: number;
  pageSize?: number;
  reportId?: string;
}): Promise<PaginatedResponse<TriageDecision>> {
  const searchParams = new URLSearchParams();
  if (params?.page) searchParams.set('page', params.page.toString());
  if (params?.pageSize) searchParams.set('pageSize', params.pageSize.toString());
  if (params?.reportId) searchParams.set('reportId', params.reportId);
  
  const query = searchParams.toString();
  return fetchApi(`/v1/triage-decisions${query ? `?${query}` : ''}`);
}

// Alerts
export async function createAlert(data: CreateAlertRequest): Promise<Alert> {
  return fetchApi('/v1/alerts', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getAlerts(params?: {
  page?: number;
  pageSize?: number;
  status?: string;
}): Promise<PaginatedResponse<Alert>> {
  const searchParams = new URLSearchParams();
  if (params?.page) searchParams.set('page', params.page.toString());
  if (params?.pageSize) searchParams.set('pageSize', params.pageSize.toString());
  if (params?.status) searchParams.set('status', params.status);
  
  const query = searchParams.toString();
  return fetchApi(`/v1/alerts${query ? `?${query}` : ''}`);
}

export async function getAlert(id: string): Promise<Alert> {
  return fetchApi(`/v1/alerts/${id}`);
}

export async function updateAlert(
  id: string,
  data: Partial<Alert>
): Promise<Alert> {
  return fetchApi(`/v1/alerts/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

// Training
export async function createTrainingEvent(data: {
  title: string;
  eventDate: string;
  location: string;
  audience?: string;
  attendanceCount?: number;
  preAvg?: number;
  postAvg?: number;
  notes?: string;
}): Promise<TrainingEvent> {
  return fetchApi('/v1/training-events', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getTrainingEvents(params?: {
  page?: number;
  pageSize?: number;
  from?: string;
  to?: string;
}): Promise<PaginatedResponse<TrainingEvent>> {
  const searchParams = new URLSearchParams();
  if (params?.page) searchParams.set('page', params.page.toString());
  if (params?.pageSize) searchParams.set('pageSize', params.pageSize.toString());
  if (params?.from) searchParams.set('from', params.from);
  if (params?.to) searchParams.set('to', params.to);
  
  const query = searchParams.toString();
  return fetchApi(`/v1/training-events${query ? `?${query}` : ''}`);
}

export async function getTrainingEvent(id: string): Promise<TrainingEvent> {
  return fetchApi(`/v1/training-events/${id}`);
}

export async function recordQuizResult(
  eventId: string,
  data: {
    participantHash?: string;
    quizType: 'pre' | 'post';
    score: number;
    maxScore?: number;
    answers?: Record<string, unknown>;
  }
): Promise<QuizResult> {
  return fetchApi(`/v1/training-events/${eventId}/results`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getTrainingStats(): Promise<TrainingStats> {
  return fetchApi('/v1/training-events/stats');
}

// Metrics
export async function getKPIMetrics(): Promise<KPIMetrics> {
  return fetchApi('/v1/metrics/kpi');
}

export async function getDashboardStats(): Promise<DashboardStats> {
  return fetchApi('/v1/metrics/dashboard');
}

// Auth
export async function login(
  email: string,
  password: string
): Promise<{ accessToken: string; tokenType: string; expiresIn: number }> {
  const result = await fetchApi<{
    accessToken: string;
    tokenType: string;
    expiresIn: number;
  }>('/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
  
  if (typeof window !== 'undefined') {
    localStorage.setItem('token', result.accessToken);
  }
  
  return result;
}

export async function getCurrentUser(): Promise<{
  id: string;
  displayName?: string;
  role: string;
}> {
  return fetchApi('/v1/auth/me');
}

export function logout(): void {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('token');
  }
}

export { ApiError };
