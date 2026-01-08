// Report types
export interface Report {
  id: string;
  category: ReportCategory;
  severitySuggested?: SeverityLevel;
  areaHint: string;
  timeWindow?: string;
  description: string;
  evidence?: string[];
  status: ReportStatus;
  createdAt: string;
  updatedAt: string;
}

export interface ReportDetail extends Report {
  triageDecisions?: TriageDecision[];
}

export type ReportCategory =
  | 'suspicious_item'
  | 'suspicious_person'
  | 'harassment_stalking'
  | 'scam_phishing'
  | 'misinformation_panic'
  | 'crowd_disorder'
  | 'infrastructure_hazard'
  | 'other';

export type ReportStatus =
  | 'submitted'
  | 'under_review'
  | 'triaged'
  | 'escalated'
  | 'closed'
  | 'spam';

export type SeverityLevel = 'S0' | 'S1' | 'S2' | 'S3' | 'S4';
export type EvidenceLevel = 'E0' | 'E1' | 'E2' | 'E3';

// Triage types
export interface TriageDecision {
  id: string;
  reportId: string;
  decision: TriageDecisionType;
  severityFinal: SeverityLevel;
  evidenceLevel?: EvidenceLevel;
  rationale?: string;
  decidedAt: string;
  decidedBy?: UserSummary;
}

export type TriageDecisionType =
  | 'accept'
  | 'reject'
  | 'needs_more_info'
  | 'escalate';

// Alert types
export interface Alert {
  id: string;
  reportId?: string;
  status: AlertStatus;
  event: string;
  urgency: CAPUrgency;
  severity: CAPSeverity;
  certainty: CAPCertainty;
  area: string;
  instruction: string;
  publicMessage?: string;
  channels?: string[];
  capXml?: string;
  createdAt: string;
  publishedAt?: string;
  updatedAt: string;
  approvedBy?: UserSummary;
}

export type AlertStatus = 'draft' | 'approved' | 'published' | 'withdrawn';
export type CAPUrgency = 'Immediate' | 'Expected' | 'Future' | 'Past' | 'Unknown';
export type CAPSeverity = 'Extreme' | 'Severe' | 'Moderate' | 'Minor' | 'Unknown';
export type CAPCertainty = 'Observed' | 'Likely' | 'Possible' | 'Unlikely' | 'Unknown';

// Training types
export interface TrainingEvent {
  id: string;
  title: string;
  eventDate: string;
  location: string;
  audience?: string;
  attendanceCount: number;
  preAvg?: number;
  postAvg?: number;
  improvement?: number;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface QuizResult {
  id: string;
  eventId?: string;
  quizType: 'pre' | 'post';
  score: number;
  maxScore: number;
  percentage: number;
  createdAt: string;
}

export interface TrainingStats {
  totalEvents: number;
  totalParticipants: number;
  averageImprovement: number;
  targetMet: boolean;
}

// Metrics types
export interface KPIMetrics {
  education: EducationKPI;
  system: SystemKPI;
  adoption: AdoptionKPI;
  governance: GovernanceKPI;
}

export interface EducationKPI {
  workshopsCount: number;
  workshopsTarget: number;
  participantsTrained: number;
  participantsTarget: number;
  prePostImprovement: number;
  improvementTarget: number;
}

export interface SystemKPI {
  medianReportToTriage: number;
  triageTimeTarget: number;
  verifiedRatio: number;
  verifiedRatioTarget: number;
  abuseRate: number;
  abuseRateTarget: number;
  alertPublishLatency: number;
  publishLatencyTarget: number;
}

export interface AdoptionKPI {
  partnerOrgs: number;
  partnerOrgsTarget: number;
  externalAdoption: number;
  externalAdoptionTarget: number;
}

export interface GovernanceKPI {
  certifiedTriagers: number;
  triagersTarget: number;
}

export interface DashboardStats {
  totalReports: number;
  reportsThisWeek: number;
  activeAlerts: number;
  recentAlerts: AlertSummary[];
  categoryBreakdown: CategoryCount[];
}

export interface AlertSummary {
  id: string;
  event: string;
  severity: CAPSeverity;
  area: string;
  createdAt: string;
}

export interface CategoryCount {
  category: ReportCategory;
  count: number;
}

// Common types
export interface UserSummary {
  id: string;
  displayName?: string;
  role: string;
}

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: Pagination;
}

export interface ErrorResponse {
  code: string;
  message: string;
  details?: Record<string, string>;
}

// Form types
export interface CreateReportRequest {
  category: ReportCategory;
  severitySuggested?: SeverityLevel;
  areaHint: string;
  timeWindow?: string;
  description: string;
  evidence?: string[];
  reporterContact?: string;
}

export interface TriageRequest {
  decision: TriageDecisionType;
  severityFinal: SeverityLevel;
  evidenceLevel?: EvidenceLevel;
  rationale?: string;
}

export interface CreateAlertRequest {
  reportId?: string;
  event: string;
  urgency: CAPUrgency;
  severity: CAPSeverity;
  certainty: CAPCertainty;
  area: string;
  instruction: string;
  publicMessage?: string;
  channels?: string[];
}
