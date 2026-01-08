package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/config"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
)

// DB wraps database connections
type DB struct {
	Gorm  *gorm.DB
	Redis *redis.Client
}

// NewDB creates a new database connection
func NewDB(cfg *config.Config) (*DB, error) {
	// Configure GORM logger
	logLevel := logger.Info
	if cfg.IsProd() {
		logLevel = logger.Warn
	}

	// Connect to PostgreSQL
	gormDB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		// Redis is optional in dev mode
		if !cfg.IsDev() {
			return nil, fmt.Errorf("failed to connect to redis: %w", err)
		}
		rdb = nil
	}

	return &DB{
		Gorm:  gormDB,
		Redis: rdb,
	}, nil
}

// AutoMigrate runs auto migration for all models (dev only)
func (db *DB) AutoMigrate() error {
	return db.Gorm.AutoMigrate(
		&model.User{},
		&model.Report{},
		&model.TriageDecision{},
		&model.Alert{},
		&model.TrainingEvent{},
		&model.TrainingParticipant{},
		&model.QuizResult{},
		&model.AuditLog{},
		&model.APIKey{},
	)
}

// Close closes all database connections
func (db *DB) Close() error {
	if db.Redis != nil {
		if err := db.Redis.Close(); err != nil {
			return err
		}
	}
	sqlDB, err := db.Gorm.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck checks database connections
func (db *DB) HealthCheck(ctx context.Context) (dbOK, redisOK bool) {
	// Check PostgreSQL
	sqlDB, err := db.Gorm.DB()
	if err == nil {
		if err := sqlDB.PingContext(ctx); err == nil {
			dbOK = true
		}
	}

	// Check Redis
	if db.Redis != nil {
		if err := db.Redis.Ping(ctx).Err(); err == nil {
			redisOK = true
		}
	}

	return
}
