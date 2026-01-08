package service

import (
	"context"
	"time"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/repository"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

// MetricsService handles KPI and metrics business logic
type MetricsService struct {
	reportRepo   *repository.ReportRepository
	triageRepo   *repository.TriageRepository
	alertRepo    *repository.AlertRepository
	trainingRepo *repository.TrainingRepository
	userRepo     *repository.UserRepository
}

// NewMetricsService creates a new metrics service
func NewMetricsService(
	reportRepo *repository.ReportRepository,
	triageRepo *repository.TriageRepository,
	alertRepo *repository.AlertRepository,
	trainingRepo *repository.TrainingRepository,
	userRepo *repository.UserRepository,
) *MetricsService {
	return &MetricsService{
		reportRepo:   reportRepo,
		triageRepo:   triageRepo,
		alertRepo:    alertRepo,
		trainingRepo: trainingRepo,
		userRepo:     userRepo,
	}
}

// GetKPIMetrics retrieves all KPI metrics
func (s *MetricsService) GetKPIMetrics(ctx context.Context) (*vo.KPIMetricsVO, error) {
	// Get training stats
	trainingStats, err := s.trainingRepo.GetTrainingStats(ctx)
	if err != nil {
		return nil, err
	}

	// Get triage metrics
	medianTriageTime, err := s.triageRepo.GetMedianTriageTime(ctx)
	if err != nil {
		medianTriageTime = 0
	}

	verifiedRatio, err := s.triageRepo.GetVerifiedRatio(ctx)
	if err != nil {
		verifiedRatio = 0
	}

	abuseRate, err := s.triageRepo.GetAbuseRate(ctx)
	if err != nil {
		abuseRate = 0
	}

	// Get alert metrics
	publishLatency, err := s.alertRepo.GetPublishLatency(ctx)
	if err != nil {
		publishLatency = 0
	}

	// Get certified triagers count
	certifiedTriagers, err := s.userRepo.CountByRole(ctx, "triager")
	if err != nil {
		certifiedTriagers = 0
	}

	return &vo.KPIMetricsVO{
		Education: vo.EducationKPIVO{
			WorkshopsCount:      int(trainingStats.TotalEvents),
			WorkshopsTarget:     12,
			ParticipantsTrained: int(trainingStats.TotalParticipants),
			ParticipantsTarget:  300,
			PrePostImprovement:  trainingStats.AverageImprovement,
			ImprovementTarget:   25.0,
		},
		System: vo.SystemKPIVO{
			MedianReportToTriage: medianTriageTime.Minutes(),
			TriageTimeTarget:     30.0,
			VerifiedRatio:        verifiedRatio,
			VerifiedRatioTarget:  60.0,
			AbuseRate:            abuseRate,
			AbuseRateTarget:      5.0,
			AlertPublishLatency:  publishLatency.Minutes(),
			PublishLatencyTarget: 15.0,
		},
		Adoption: vo.AdoptionKPIVO{
			PartnerOrgs:            0, // Manually tracked
			PartnerOrgsTarget:      4,
			ExternalAdoption:       0, // Manually tracked
			ExternalAdoptionTarget: 2,
		},
		Governance: vo.GovernanceKPIVO{
			CertifiedTriagers: int(certifiedTriagers),
			TriagersTarget:    15,
		},
	}, nil
}

// GetDashboardStats retrieves public dashboard statistics
func (s *MetricsService) GetDashboardStats(ctx context.Context) (*vo.DashboardStatsVO, error) {
	// Get report stats
	reportStats, err := s.reportRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	// Get reports this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	reportsThisWeek := int64(0)
	for status, count := range reportStats.ByStatus {
		if status != "spam" {
			reportsThisWeek += count
		}
	}

	// Get active alerts
	activeAlerts, err := s.alertRepo.CountByStatus(ctx, "published")
	if err != nil {
		activeAlerts = 0
	}

	// Get recent alerts
	recentAlerts, err := s.alertRepo.GetRecentAlerts(ctx, 5)
	if err != nil {
		recentAlerts = nil
	}

	alertSummaries := make([]vo.AlertSummaryVO, len(recentAlerts))
	for i, a := range recentAlerts {
		alertSummaries[i] = vo.AlertSummaryVO{
			ID:        a.ID.String(),
			Event:     a.Event,
			Severity:  a.Severity,
			Area:      a.Area,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		}
	}

	// Get category breakdown
	categoryBreakdown := make([]vo.CategoryCountVO, 0, len(reportStats.ByCategory))
	for category, count := range reportStats.ByCategory {
		categoryBreakdown = append(categoryBreakdown, vo.CategoryCountVO{
			Category: category,
			Count:    int(count),
		})
	}

	return &vo.DashboardStatsVO{
		TotalReports:      int(reportStats.Total),
		ReportsThisWeek:   int(reportsThisWeek),
		ActiveAlerts:      int(activeAlerts),
		RecentAlerts:      alertSummaries,
		CategoryBreakdown: categoryBreakdown,
	}, nil
}
