package types

import (
  "encoding/xml"
  "time"

  "github.com/google/uuid"
)

type CreateReportRequest struct {
  Category    string   `json:"category" binding:"required"`
  AreaHint    string   `json:"areaHint" binding:"required"`
  TimeWindow  string   `json:"timeWindow"`
  Description string   `json:"description" binding:"required"`
  Evidence    []string `json:"evidence"`
}

type Report struct {
  ID          string    `json:"id"`
  CreatedAt   time.Time `json:"createdAt"`
  Category    string    `json:"category"`
  AreaHint    string    `json:"areaHint"`
  TimeWindow  string    `json:"timeWindow,omitempty"`
  Description string    `json:"description"`
  Evidence    []string  `json:"evidence,omitempty"`
  Status      string    `json:"status"`
}

type TriageRequest struct {
  Decision     string `json:"decision" binding:"required"`       // accept/reject/needs_more_info/escalate
  SeverityFinal string `json:"severityFinal" binding:"required"` // S0..S4
  EvidenceLevel string `json:"evidenceLevel"`                    // E0..E3
  Rationale     string `json:"rationale"`
}

type TriageDecision struct {
  ID          string    `json:"id"`
  ReportID    string    `json:"reportId"`
  DecidedAt   time.Time `json:"decidedAt"`
  Decision    string    `json:"decision"`
  SeverityFinal string  `json:"severityFinal"`
  EvidenceLevel string  `json:"evidenceLevel,omitempty"`
  Rationale   string    `json:"rationale,omitempty"`
}

type CreateAlertRequest struct {
  Event       string `json:"event" binding:"required"`
  Urgency     string `json:"urgency" binding:"required"`
  Severity    string `json:"severity" binding:"required"`
  Certainty   string `json:"certainty" binding:"required"`
  Area        string `json:"area" binding:"required"`
  Instruction string `json:"instruction" binding:"required"`
}

type Alert struct {
  ID          string    `json:"id"`
  CreatedAt   time.Time `json:"createdAt"`
  Status      string    `json:"status"`
  Event       string    `json:"event"`
  Urgency     string    `json:"urgency"`
  Severity    string    `json:"severity"`
  Certainty   string    `json:"certainty"`
  Area        string    `json:"area"`
  Instruction string    `json:"instruction"`
  CAPXML      string    `json:"capXml"`
}

func NewID() string { return uuid.NewString() }

type capAlert struct {
  XMLName    xml.Name `xml:"alert"`
  Xmlns      string   `xml:"xmlns,attr,omitempty"`
  Identifier string   `xml:"identifier"`
  Sender     string   `xml:"sender"`
  Sent       string   `xml:"sent"`
  Status     string   `xml:"status"`
  MsgType    string   `xml:"msgType"`
  Scope      string   `xml:"scope"`
  Info       capInfo  `xml:"info"`
}

type capInfo struct {
  Category    string `xml:"category"`
  Event       string `xml:"event"`
  Urgency     string `xml:"urgency"`
  Severity    string `xml:"severity"`
  Certainty   string `xml:"certainty"`
  Instruction string `xml:"instruction"`
  Area        capArea `xml:"area"`
}

type capArea struct {
  AreaDesc string `xml:"areaDesc"`
}

func BuildCAPXML(req CreateAlertRequest) string {
  a := capAlert{
    Xmlns: "urn:oasis:names:tc:emergency:cap:1.2",
    Identifier: NewID(),
    Sender: "the-hive@example.invalid",
    Sent: time.Now().UTC().Format(time.RFC3339),
    Status: "Actual",
    MsgType: "Alert",
    Scope: "Public",
    Info: capInfo{
      Category: "Safety",
      Event: req.Event,
      Urgency: req.Urgency,
      Severity: req.Severity,
      Certainty: req.Certainty,
      Instruction: req.Instruction,
      Area: capArea{AreaDesc: req.Area},
    },
  }
  b, err := xml.MarshalIndent(a, "", "  ")
  if err != nil {
    return ""
  }
  return xml.Header + string(b)
}
