package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/dto"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/repository"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

var (
	ErrReportNotFound = errors.New("report not found")
)

// ReportService handles report business logic
type ReportService struct {
	reportRepo *repository.ReportRepository
	auditRepo  *repository.AuditRepository
}

// NewReportService creates a new report service
func NewReportService(reportRepo *repository.ReportRepository, auditRepo *repository.AuditRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		auditRepo:  auditRepo,
	}
}

// Create creates a new report
func (s *ReportService) Create(ctx context.Context, req dto.CreateReportRequest, actorIP string) (*vo.ReportVO, error) {
	report := &model.Report{
		Category:           req.Category,
		SeveritySuggested:  req.SeveritySuggested,
		AreaHint:           req.AreaHint,
		TimeWindow:         req.TimeWindow,
		Description:        req.Description,
		EvidenceRefs:       req.Evidence,
		ReporterContactRef: req.ReporterContact,
		Status:             model.StatusSubmitted,
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, err
	}

	// Create audit log
	s.auditRepo.Create(ctx, &model.AuditLog{
		ActorIP:    actorIP,
		Action:     model.ActionCreate,
		ObjectType: model.ObjectTypeReport,
		ObjectID:   &report.ID,
		Diff: model.JSONMap{
			"category": report.Category,
			"status":   report.Status,
		},
	})

	return s.toReportVO(report), nil
}

// GetByID retrieves a report by ID
func (s *ReportService) GetByID(ctx context.Context, id string) (*vo.ReportDetailVO, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrReportNotFound
	}

	report, err := s.reportRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, ErrReportNotFound
	}

	return s.toReportDetailVO(report), nil
}

// List retrieves reports with pagination
func (s *ReportService) List(ctx context.Context, query dto.ListReportsQuery) (*vo.ReportListVO, error) {
	params := repository.ListReportParams{
		Page:     query.Page,
		PageSize: query.PageSize,
		Status:   query.Status,
		Category: query.Category,
		SortBy:   toSnakeCase(query.SortBy),
		SortDir:  query.SortDir,
	}

	reports, total, err := s.reportRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	reportVOs := make([]vo.ReportVO, len(reports))
	for i, r := range reports {
		reportVOs[i] = *s.toReportVO(&r)
	}

	return &vo.ReportListVO{
		Data: reportVOs,
		Pagination: vo.PaginationVO{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      total,
			TotalPages: int((total + int64(query.PageSize) - 1) / int64(query.PageSize)),
		},
	}, nil
}

// toReportVO converts a report model to VO
func (s *ReportService) toReportVO(report *model.Report) *vo.ReportVO {
	return &vo.ReportVO{
		ID:                report.ID.String(),
		Category:          report.Category,
		SeveritySuggested: report.SeveritySuggested,
		AreaHint:          report.AreaHint,
		TimeWindow:        report.TimeWindow,
		Description:       report.Description,
		Evidence:          report.EvidenceRefs,
		Status:            report.Status,
		CreatedAt:         report.CreatedAt,
		UpdatedAt:         report.UpdatedAt,
	}
}

// toReportDetailVO converts a report model to detailed VO with triage decisions
func (s *ReportService) toReportDetailVO(report *model.Report) *vo.ReportDetailVO {
	detail := &vo.ReportDetailVO{
		ReportVO: *s.toReportVO(report),
	}

	if len(report.TriageDecisions) > 0 {
		detail.TriageDecisions = make([]vo.TriageDecisionVO, len(report.TriageDecisions))
		for i, td := range report.TriageDecisions {
			copier.Copy(&detail.TriageDecisions[i], &td)
			detail.TriageDecisions[i].ID = td.ID.String()
			detail.TriageDecisions[i].ReportID = td.ReportID.String()
		}
	}

	return detail
}

// toSnakeCase converts camelCase to snake_case for DB column names
func toSnakeCase(s string) string {
	switch s {
	case "createdAt":
		return "created_at"
	case "updatedAt":
		return "updated_at"
	default:
		return s
	}
}
