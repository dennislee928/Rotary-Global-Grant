package httpapi

import (
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/google/uuid"

  "github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/store"
  "github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/types"
)

func NewRouter() *gin.Engine {
  r := gin.New()
  r.Use(gin.Logger(), gin.Recovery())

  mem := store.NewMemoryStore()

  r.GET("/healthz", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok", "ts": time.Now().UTC()})
  })

  v1 := r.Group("/v1")

  v1.POST("/reports", func(c *gin.Context) {
    var req types.CreateReportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    rep := types.Report{
      ID:        uuid.NewString(),
      CreatedAt: time.Now().UTC(),
      Category:  req.Category,
      AreaHint:  req.AreaHint,
      TimeWindow:req.TimeWindow,
      Description: req.Description,
      Status:    "submitted",
      Evidence:  req.Evidence,
    }
    mem.CreateReport(rep)
    c.JSON(http.StatusCreated, rep)
  })

  v1.GET("/reports", func(c *gin.Context) {
    c.JSON(http.StatusOK, mem.ListReports())
  })

  v1.POST("/reports/:id/triage", func(c *gin.Context) {
    id := c.Param("id")
    var req types.TriageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    decision, err := mem.TriageReport(id, req)
    if err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
      return
    }
    c.JSON(http.StatusOK, decision)
  })

  v1.POST("/alerts", func(c *gin.Context) {
    var req types.CreateAlertRequest
    if err := c.ShouldBindJSON(&req); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    alert := types.Alert{
      ID: uuid.NewString(),
      CreatedAt: time.Now().UTC(),
      Status: "draft",
      Event: req.Event,
      Urgency: req.Urgency,
      Severity: req.Severity,
      Certainty: req.Certainty,
      Area: req.Area,
      Instruction: req.Instruction,
      CAPXML: types.BuildCAPXML(req),
    }
    mem.CreateAlert(alert)
    c.JSON(http.StatusCreated, alert)
  })

  return r
}
