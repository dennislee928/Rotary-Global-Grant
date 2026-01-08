package store

import (
  "errors"
  "sync"
  "time"

  "github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/types"
)

type MemoryStore struct {
  mu sync.Mutex
  reports map[string]types.Report
  triage  map[string]types.TriageDecision
  alerts  map[string]types.Alert
}

func NewMemoryStore() *MemoryStore {
  return &MemoryStore{
    reports: map[string]types.Report{},
    triage:  map[string]types.TriageDecision{},
    alerts:  map[string]types.Alert{},
  }
}

func (s *MemoryStore) CreateReport(r types.Report) {
  s.mu.Lock()
  defer s.mu.Unlock()
  s.reports[r.ID] = r
}

func (s *MemoryStore) ListReports() []types.Report {
  s.mu.Lock()
  defer s.mu.Unlock()
  out := make([]types.Report, 0, len(s.reports))
  for _, v := range s.reports {
    out = append(out, v)
  }
  return out
}

func (s *MemoryStore) TriageReport(reportID string, req types.TriageRequest) (types.TriageDecision, error) {
  s.mu.Lock()
  defer s.mu.Unlock()

  rep, ok := s.reports[reportID]
  if !ok {
    return types.TriageDecision{}, errors.New("report not found")
  }
  rep.Status = "triaged"
  s.reports[reportID] = rep

  d := types.TriageDecision{
    ID: types.NewID(),
    ReportID: reportID,
    DecidedAt: time.Now().UTC(),
    Decision: req.Decision,
    SeverityFinal: req.SeverityFinal,
    EvidenceLevel: req.EvidenceLevel,
    Rationale: req.Rationale,
  }
  s.triage[d.ID] = d
  return d, nil
}

func (s *MemoryStore) CreateAlert(a types.Alert) {
  s.mu.Lock()
  defer s.mu.Unlock()
  s.alerts[a.ID] = a
}
