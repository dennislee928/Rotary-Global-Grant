package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/config"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/handler"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/middleware"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/repository"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/service"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

// Server holds all dependencies for the HTTP server
type Server struct {
	Router *gin.Engine
	DB     *repository.DB
	Config *config.Config
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config) (*Server, error) {
	// Set Gin mode
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := repository.NewDB(cfg)
	if err != nil {
		return nil, err
	}

	// Create router
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	reportRepo := repository.NewReportRepository(db)
	triageRepo := repository.NewTriageRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	trainingRepo := repository.NewTrainingRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// Create services
	authSvc := service.NewAuthService(userRepo, auditRepo, cfg.JWTSecret, cfg.JWTExpiration)
	reportSvc := service.NewReportService(reportRepo, auditRepo)
	triageSvc := service.NewTriageService(triageRepo, reportRepo, auditRepo)
	alertSvc := service.NewAlertService(alertRepo, auditRepo, cfg.CAPSender)
	trainingSvc := service.NewTrainingService(trainingRepo, auditRepo)
	metricsSvc := service.NewMetricsService(reportRepo, triageRepo, alertRepo, trainingRepo, userRepo)

	// Create handlers
	authHandler := handler.NewAuthHandler(authSvc)
	reportHandler := handler.NewReportHandler(reportSvc)
	triageHandler := handler.NewTriageHandler(triageSvc)
	alertHandler := handler.NewAlertHandler(alertSvc)
	trainingHandler := handler.NewTrainingHandler(trainingSvc)
	metricsHandler := handler.NewMetricsHandler(metricsSvc)

	// Rate limiter
	var rateLimiter gin.HandlerFunc
	if db.Redis != nil {
		rl := middleware.NewRateLimiter(db.Redis, cfg.RateLimitRequests, cfg.RateLimitWindow)
		rateLimiter = rl.Middleware()
	} else {
		rl := middleware.NewInMemoryRateLimiter(cfg.RateLimitRequests)
		rateLimiter = rl.Middleware()
	}

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		dbOK, redisOK := db.HealthCheck(context.Background())
		dbStatus := "connected"
		if !dbOK {
			dbStatus = "disconnected"
		}
		redisStatus := "connected"
		if !redisOK {
			redisStatus = "disconnected"
		}

		c.JSON(http.StatusOK, vo.HealthVO{
			Status:    "ok",
			Timestamp: time.Now().UTC(),
			Version:   "0.1.0",
			Database:  dbStatus,
			Redis:     redisStatus,
		})
	})

	// API v1 routes
	v1 := r.Group("/v1")
	v1.Use(rateLimiter)

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Protected auth routes
	authProtected := v1.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware(authSvc))
	{
		authProtected.GET("/me", authHandler.Me)
		authProtected.POST("/api-keys", authHandler.CreateAPIKey)
	}

	// Reports routes
	reports := v1.Group("/reports")
	{
		// Public: create report
		reports.POST("", reportHandler.Create)

		// Protected: list and get reports
		reportsProtected := reports.Group("")
		reportsProtected.Use(middleware.AuthMiddleware(authSvc))
		{
			reportsProtected.GET("", reportHandler.List)
			reportsProtected.GET("/:id", reportHandler.GetByID)
			reportsProtected.POST("/:id/triage", 
				middleware.RoleMiddleware(model.RoleAdmin, model.RoleTriager),
				triageHandler.TriageReport,
			)
		}
	}

	// Triage decisions routes (protected)
	triage := v1.Group("/triage-decisions")
	triage.Use(middleware.AuthMiddleware(authSvc))
	{
		triage.GET("", triageHandler.List)
	}

	// Alerts routes
	alerts := v1.Group("/alerts")
	alerts.Use(middleware.AuthMiddleware(authSvc))
	{
		alerts.POST("", 
			middleware.RoleMiddleware(model.RoleAdmin, model.RoleTriager),
			alertHandler.Create,
		)
		alerts.GET("", alertHandler.List)
		alerts.GET("/:id", alertHandler.GetByID)
		alerts.PATCH("/:id", 
			middleware.RoleMiddleware(model.RoleAdmin),
			alertHandler.Update,
		)
	}

	// Training events routes
	training := v1.Group("/training-events")
	{
		// Public: get training stats
		training.GET("/stats", trainingHandler.GetStats)

		// Protected routes
		trainingProtected := training.Group("")
		trainingProtected.Use(middleware.AuthMiddleware(authSvc))
		{
			trainingProtected.POST("", 
				middleware.RoleMiddleware(model.RoleAdmin, model.RoleEducator),
				trainingHandler.Create,
			)
			trainingProtected.GET("", trainingHandler.List)
			trainingProtected.GET("/:id", trainingHandler.GetByID)
			trainingProtected.POST("/:id/results", trainingHandler.RecordQuizResult)
		}
	}

	// Metrics routes
	metrics := v1.Group("/metrics")
	{
		// Public: dashboard stats
		metrics.GET("/dashboard", metricsHandler.GetDashboardStats)

		// Protected: KPI metrics
		metricsProtected := metrics.Group("")
		metricsProtected.Use(middleware.AuthMiddleware(authSvc))
		{
			metricsProtected.GET("/kpi", 
				middleware.RoleMiddleware(model.RoleAdmin, model.RoleAuditor),
				metricsHandler.GetKPIMetrics,
			)
		}
	}

	return &Server{
		Router: r,
		DB:     db,
		Config: cfg,
	}, nil
}

// Run starts the HTTP server
func (s *Server) Run() error {
	return s.Router.Run(s.Config.HTTPAddr)
}

// Close cleans up server resources
func (s *Server) Close() error {
	return s.DB.Close()
}
