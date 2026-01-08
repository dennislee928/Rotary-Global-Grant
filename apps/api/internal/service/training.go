package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/dto"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/repository"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

var (
	ErrTrainingEventNotFound = errors.New("training event not found")
)

// TrainingService handles training event business logic
type TrainingService struct {
	trainingRepo *repository.TrainingRepository
	auditRepo    *repository.AuditRepository
}

// NewTrainingService creates a new training service
func NewTrainingService(trainingRepo *repository.TrainingRepository, auditRepo *repository.AuditRepository) *TrainingService {
	return &TrainingService{
		trainingRepo: trainingRepo,
		auditRepo:    auditRepo,
	}
}

// CreateEvent creates a new training event
func (s *TrainingService) CreateEvent(ctx context.Context, req dto.CreateTrainingEventRequest, userID *uuid.UUID, actorIP string) (*vo.TrainingEventVO, error) {
	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	var preAvg, postAvg *float64
	if req.PreAvg > 0 {
		preAvg = &req.PreAvg
	}
	if req.PostAvg > 0 {
		postAvg = &req.PostAvg
	}

	event := &model.TrainingEvent{
		Title:           req.Title,
		EventDate:       eventDate,
		Location:        req.Location,
		Audience:        req.Audience,
		AttendanceCount: req.AttendanceCount,
		PreAvg:          preAvg,
		PostAvg:         postAvg,
		Notes:           req.Notes,
		CreatedBy:       userID,
	}

	if err := s.trainingRepo.CreateEvent(ctx, event); err != nil {
		return nil, err
	}

	// Create audit log
	s.auditRepo.Create(ctx, &model.AuditLog{
		ActorID:    userID,
		ActorIP:    actorIP,
		Action:     model.ActionCreate,
		ObjectType: model.ObjectTypeTraining,
		ObjectID:   &event.ID,
		Diff: model.JSONMap{
			"title":    event.Title,
			"date":     req.EventDate,
			"location": event.Location,
		},
	})

	return s.toTrainingEventVO(event), nil
}

// GetEventByID retrieves a training event by ID
func (s *TrainingService) GetEventByID(ctx context.Context, id string) (*vo.TrainingEventVO, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrTrainingEventNotFound
	}

	event, err := s.trainingRepo.GetEventByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrTrainingEventNotFound
	}

	return s.toTrainingEventVO(event), nil
}

// ListEvents retrieves training events with pagination
func (s *TrainingService) ListEvents(ctx context.Context, query dto.ListTrainingEventsQuery) (*vo.TrainingEventListVO, error) {
	var from, to time.Time
	if query.From != "" {
		var err error
		from, err = time.Parse("2006-01-02", query.From)
		if err != nil {
			return nil, errors.New("invalid from date format")
		}
	}
	if query.To != "" {
		var err error
		to, err = time.Parse("2006-01-02", query.To)
		if err != nil {
			return nil, errors.New("invalid to date format")
		}
	}

	params := repository.ListTrainingParams{
		Page:     query.Page,
		PageSize: query.PageSize,
		From:     from,
		To:       to,
	}

	events, total, err := s.trainingRepo.ListEvents(ctx, params)
	if err != nil {
		return nil, err
	}

	eventVOs := make([]vo.TrainingEventVO, len(events))
	for i, e := range events {
		eventVOs[i] = *s.toTrainingEventVO(&e)
	}

	return &vo.TrainingEventListVO{
		Data: eventVOs,
		Pagination: vo.PaginationVO{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      total,
			TotalPages: int((total + int64(query.PageSize) - 1) / int64(query.PageSize)),
		},
	}, nil
}

// RecordQuizResult records a quiz result for a training event
func (s *TrainingService) RecordQuizResult(ctx context.Context, eventID string, req dto.RecordQuizResultRequest) (*vo.QuizResultVO, error) {
	uid, err := uuid.Parse(eventID)
	if err != nil {
		return nil, ErrTrainingEventNotFound
	}

	// Check if event exists
	event, err := s.trainingRepo.GetEventByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrTrainingEventNotFound
	}

	maxScore := req.MaxScore
	if maxScore == 0 {
		maxScore = 100
	}

	result := &model.QuizResult{
		EventID:         &uid,
		ParticipantHash: req.ParticipantHash,
		QuizType:        req.QuizType,
		Score:           req.Score,
		MaxScore:        maxScore,
		Answers:         req.Answers,
	}

	if err := s.trainingRepo.CreateQuizResult(ctx, result); err != nil {
		return nil, err
	}

	return s.toQuizResultVO(result), nil
}

// GetStats retrieves training statistics
func (s *TrainingService) GetStats(ctx context.Context) (*vo.TrainingStatsVO, error) {
	stats, err := s.trainingRepo.GetTrainingStats(ctx)
	if err != nil {
		return nil, err
	}

	return &vo.TrainingStatsVO{
		TotalEvents:        int(stats.TotalEvents),
		TotalParticipants:  int(stats.TotalParticipants),
		AverageImprovement: stats.AverageImprovement,
		TargetMet:          stats.TotalParticipants >= 300 && stats.AverageImprovement >= 25,
	}, nil
}

// toTrainingEventVO converts a training event model to VO
func (s *TrainingService) toTrainingEventVO(event *model.TrainingEvent) *vo.TrainingEventVO {
	result := &vo.TrainingEventVO{
		ID:              event.ID.String(),
		Title:           event.Title,
		EventDate:       event.EventDate.Format("2006-01-02"),
		Location:        event.Location,
		Audience:        event.Audience,
		AttendanceCount: event.AttendanceCount,
		PreAvg:          event.PreAvg,
		PostAvg:         event.PostAvg,
		Notes:           event.Notes,
		CreatedAt:       event.CreatedAt,
		UpdatedAt:       event.UpdatedAt,
	}

	// Calculate improvement
	if event.PreAvg != nil && event.PostAvg != nil {
		improvement := *event.PostAvg - *event.PreAvg
		result.Improvement = &improvement
	}

	return result
}

// toQuizResultVO converts a quiz result model to VO
func (s *TrainingService) toQuizResultVO(result *model.QuizResult) *vo.QuizResultVO {
	vo := &vo.QuizResultVO{
		ID:         result.ID.String(),
		QuizType:   result.QuizType,
		Score:      result.Score,
		MaxScore:   result.MaxScore,
		Percentage: (result.Score / result.MaxScore) * 100,
		CreatedAt:  result.CreatedAt,
	}

	if result.EventID != nil {
		vo.EventID = result.EventID.String()
	}

	return vo
}
