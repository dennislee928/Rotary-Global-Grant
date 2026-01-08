package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/dto"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/pkg/cap"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/repository"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

var (
	ErrAlertNotFound     = errors.New("alert not found")
	ErrInvalidTransition = errors.New("invalid status transition")
)

// AlertService handles alert business logic
type AlertService struct {
	alertRepo *repository.AlertRepository
	auditRepo *repository.AuditRepository
	capSender string
}

// NewAlertService creates a new alert service
func NewAlertService(alertRepo *repository.AlertRepository, auditRepo *repository.AuditRepository, capSender string) *AlertService {
	return &AlertService{
		alertRepo: alertRepo,
		auditRepo: auditRepo,
		capSender: capSender,
	}
}

// Create creates a new alert
func (s *AlertService) Create(ctx context.Context, req dto.CreateAlertRequest, userID *uuid.UUID, actorIP string) (*vo.AlertVO, error) {
	var reportID *uuid.UUID
	if req.ReportID != "" {
		uid, err := uuid.Parse(req.ReportID)
		if err != nil {
			return nil, errors.New("invalid report ID")
		}
		reportID = &uid
	}

	// Build CAP XML
	capXML := cap.BuildCAPXML(cap.CAPParams{
		Sender:      s.capSender,
		Event:       req.Event,
		Urgency:     req.Urgency,
		Severity:    req.Severity,
		Certainty:   req.Certainty,
		Area:        req.Area,
		Instruction: req.Instruction,
	})

	alert := &model.Alert{
		ReportID:      reportID,
		Status:        model.AlertStatusDraft,
		Event:         req.Event,
		Urgency:       req.Urgency,
		Severity:      req.Severity,
		Certainty:     req.Certainty,
		Area:          req.Area,
		Instruction:   req.Instruction,
		PublicMessage: req.PublicMessage,
		CAPXML:        capXML,
		Channels:      req.Channels,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, err
	}

	// Create audit log
	s.auditRepo.Create(ctx, &model.AuditLog{
		ActorID:    userID,
		ActorIP:    actorIP,
		Action:     model.ActionCreate,
		ObjectType: model.ObjectTypeAlert,
		ObjectID:   &alert.ID,
		Diff: model.JSONMap{
			"event":    alert.Event,
			"severity": alert.Severity,
			"status":   alert.Status,
		},
	})

	return s.toAlertVO(alert), nil
}

// GetByID retrieves an alert by ID
func (s *AlertService) GetByID(ctx context.Context, id string) (*vo.AlertVO, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrAlertNotFound
	}

	alert, err := s.alertRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if alert == nil {
		return nil, ErrAlertNotFound
	}

	return s.toAlertVO(alert), nil
}

// List retrieves alerts with pagination
func (s *AlertService) List(ctx context.Context, query dto.ListAlertsQuery) (*vo.AlertListVO, error) {
	params := repository.ListAlertParams{
		Page:     query.Page,
		PageSize: query.PageSize,
		Status:   query.Status,
	}

	alerts, total, err := s.alertRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	alertVOs := make([]vo.AlertVO, len(alerts))
	for i, a := range alerts {
		alertVOs[i] = *s.toAlertVO(&a)
	}

	return &vo.AlertListVO{
		Data: alertVOs,
		Pagination: vo.PaginationVO{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      total,
			TotalPages: int((total + int64(query.PageSize) - 1) / int64(query.PageSize)),
		},
	}, nil
}

// Update updates an alert
func (s *AlertService) Update(ctx context.Context, id string, req dto.UpdateAlertRequest, userID *uuid.UUID, actorIP string) (*vo.AlertVO, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrAlertNotFound
	}

	alert, err := s.alertRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if alert == nil {
		return nil, ErrAlertNotFound
	}

	// Validate status transition
	if req.Status != "" && !isValidStatusTransition(alert.Status, req.Status) {
		return nil, ErrInvalidTransition
	}

	// Update fields
	changes := make(model.JSONMap)
	if req.Status != "" {
		changes["status"] = map[string]string{"from": alert.Status, "to": req.Status}
		alert.Status = req.Status
		if req.Status == model.AlertStatusApproved {
			alert.ApprovedBy = userID
		} else if req.Status == model.AlertStatusPublished {
			now := time.Now().UTC()
			alert.PublishedAt = &now
		}
	}
	if req.Event != "" {
		alert.Event = req.Event
	}
	if req.Urgency != "" {
		alert.Urgency = req.Urgency
	}
	if req.Severity != "" {
		alert.Severity = req.Severity
	}
	if req.Certainty != "" {
		alert.Certainty = req.Certainty
	}
	if req.Area != "" {
		alert.Area = req.Area
	}
	if req.Instruction != "" {
		alert.Instruction = req.Instruction
	}
	if req.PublicMessage != "" {
		alert.PublicMessage = req.PublicMessage
	}
	if req.Channels != nil {
		alert.Channels = req.Channels
	}

	// Regenerate CAP XML if content changed
	alert.CAPXML = cap.BuildCAPXML(cap.CAPParams{
		Sender:      s.capSender,
		Event:       alert.Event,
		Urgency:     alert.Urgency,
		Severity:    alert.Severity,
		Certainty:   alert.Certainty,
		Area:        alert.Area,
		Instruction: alert.Instruction,
	})

	alert.UpdatedAt = time.Now().UTC()

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, err
	}

	// Determine audit action
	action := model.ActionUpdate
	if req.Status == model.AlertStatusApproved {
		action = model.ActionApprove
	} else if req.Status == model.AlertStatusPublished {
		action = model.ActionPublish
	} else if req.Status == model.AlertStatusWithdrawn {
		action = model.ActionWithdraw
	}

	// Create audit log
	s.auditRepo.Create(ctx, &model.AuditLog{
		ActorID:    userID,
		ActorIP:    actorIP,
		Action:     action,
		ObjectType: model.ObjectTypeAlert,
		ObjectID:   &alert.ID,
		Diff:       changes,
	})

	return s.toAlertVO(alert), nil
}

// GetActiveAlerts retrieves currently active alerts
func (s *AlertService) GetActiveAlerts(ctx context.Context, limit int) ([]vo.AlertSummaryVO, error) {
	alerts, err := s.alertRepo.GetActiveAlerts(ctx, limit)
	if err != nil {
		return nil, err
	}

	summaries := make([]vo.AlertSummaryVO, len(alerts))
	for i, a := range alerts {
		summaries[i] = vo.AlertSummaryVO{
			ID:        a.ID.String(),
			Event:     a.Event,
			Severity:  a.Severity,
			Area:      a.Area,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		}
	}

	return summaries, nil
}

// toAlertVO converts an alert model to VO
func (s *AlertService) toAlertVO(alert *model.Alert) *vo.AlertVO {
	result := &vo.AlertVO{
		ID:            alert.ID.String(),
		Status:        alert.Status,
		Event:         alert.Event,
		Urgency:       alert.Urgency,
		Severity:      alert.Severity,
		Certainty:     alert.Certainty,
		Area:          alert.Area,
		Instruction:   alert.Instruction,
		PublicMessage: alert.PublicMessage,
		Channels:      alert.Channels,
		CAPXML:        alert.CAPXML,
		CreatedAt:     alert.CreatedAt,
		PublishedAt:   alert.PublishedAt,
		UpdatedAt:     alert.UpdatedAt,
	}

	if alert.ReportID != nil {
		result.ReportID = alert.ReportID.String()
	}

	if alert.Approver != nil {
		result.ApprovedBy = &vo.UserSummaryVO{
			ID:          alert.Approver.ID.String(),
			DisplayName: alert.Approver.DisplayName,
			Role:        alert.Approver.Role,
		}
	}

	return result
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		model.AlertStatusDraft:     {model.AlertStatusApproved, model.AlertStatusWithdrawn},
		model.AlertStatusApproved:  {model.AlertStatusPublished, model.AlertStatusWithdrawn, model.AlertStatusDraft},
		model.AlertStatusPublished: {model.AlertStatusWithdrawn},
		model.AlertStatusWithdrawn: {model.AlertStatusDraft},
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
