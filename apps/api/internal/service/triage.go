package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/dto"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/repository"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

var (
	ErrTriageNotFound = errors.New("triage decision not found")
)

// TriageService handles triage business logic
type TriageService struct {
	triageRepo *repository.TriageRepository
	reportRepo *repository.ReportRepository
	auditRepo  *repository.AuditRepository
}

// NewTriageService creates a new triage service
func NewTriageService(
	triageRepo *repository.TriageRepository,
	reportRepo *repository.ReportRepository,
	auditRepo *repository.AuditRepository,
) *TriageService {
	return &TriageService{
		triageRepo: triageRepo,
		reportRepo: reportRepo,
		auditRepo:  auditRepo,
	}
}

// TriageReport creates a triage decision for a report
func (s *TriageService) TriageReport(ctx context.Context, reportID string, req dto.TriageRequest, userID *uuid.UUID, actorIP string) (*vo.TriageDecisionVO, error) {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return nil, ErrReportNotFound
	}

	// Check if report exists
	report, err := s.reportRepo.GetByID(ctx, reportUUID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, ErrReportNotFound
	}

	// Create audit hash
	auditData := map[string]interface{}{
		"reportId":      reportID,
		"decision":      req.Decision,
		"severityFinal": req.SeverityFinal,
		"evidenceLevel": req.EvidenceLevel,
		"rationale":     req.Rationale,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	}
	auditHash := generateAuditHash(auditData)

	decision := &model.TriageDecision{
		ReportID:      reportUUID,
		DecidedBy:     userID,
		Decision:      req.Decision,
		SeverityFinal: req.SeverityFinal,
		EvidenceLevel: req.EvidenceLevel,
		Rationale:     req.Rationale,
		AuditHash:     auditHash,
		DecidedAt:     time.Now().UTC(),
	}

	if err := s.triageRepo.Create(ctx, decision); err != nil {
		return nil, err
	}

	// Update report status based on decision
	newStatus := model.StatusTriaged
	if req.Decision == model.DecisionEscalate {
		newStatus = model.StatusEscalated
	} else if req.Decision == model.DecisionReject {
		newStatus = model.StatusClosed
	}
	s.reportRepo.UpdateStatus(ctx, reportUUID, newStatus)

	// Create audit log
	s.auditRepo.Create(ctx, &model.AuditLog{
		ActorID:    userID,
		ActorIP:    actorIP,
		Action:     model.ActionTriage,
		ObjectType: model.ObjectTypeReport,
		ObjectID:   &reportUUID,
		Diff: model.JSONMap{
			"decision":      req.Decision,
			"severityFinal": req.SeverityFinal,
			"auditHash":     auditHash,
		},
	})

	return s.toTriageDecisionVO(decision), nil
}

// List retrieves triage decisions with pagination
func (s *TriageService) List(ctx context.Context, query dto.ListTriageDecisionsQuery) (*vo.TriageDecisionListVO, error) {
	var reportID uuid.UUID
	if query.ReportID != "" {
		var err error
		reportID, err = uuid.Parse(query.ReportID)
		if err != nil {
			return nil, errors.New("invalid report ID")
		}
	}

	params := repository.ListTriageParams{
		Page:     query.Page,
		PageSize: query.PageSize,
		ReportID: reportID,
		Decision: query.Decision,
	}

	decisions, total, err := s.triageRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	decisionVOs := make([]vo.TriageDecisionVO, len(decisions))
	for i, d := range decisions {
		decisionVOs[i] = *s.toTriageDecisionVO(&d)
	}

	return &vo.TriageDecisionListVO{
		Data: decisionVOs,
		Pagination: vo.PaginationVO{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      total,
			TotalPages: int((total + int64(query.PageSize) - 1) / int64(query.PageSize)),
		},
	}, nil
}

// toTriageDecisionVO converts a triage decision model to VO
func (s *TriageService) toTriageDecisionVO(decision *model.TriageDecision) *vo.TriageDecisionVO {
	result := &vo.TriageDecisionVO{
		ID:            decision.ID.String(),
		ReportID:      decision.ReportID.String(),
		Decision:      decision.Decision,
		SeverityFinal: decision.SeverityFinal,
		EvidenceLevel: decision.EvidenceLevel,
		Rationale:     decision.Rationale,
		DecidedAt:     decision.DecidedAt,
	}

	if decision.Decider != nil {
		result.DecidedBy = &vo.UserSummaryVO{
			ID:          decision.Decider.ID.String(),
			DisplayName: decision.Decider.DisplayName,
			Role:        decision.Decider.Role,
		}
	}

	return result
}

// generateAuditHash creates a SHA256 hash of the audit data
func generateAuditHash(data map[string]interface{}) string {
	jsonBytes, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:])
}
