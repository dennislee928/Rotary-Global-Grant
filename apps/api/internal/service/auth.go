package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/dto"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/repository"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo      *repository.UserRepository
	auditRepo     *repository.AuditRepository
	jwtSecret     []byte
	jwtExpiration time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repository.UserRepository,
	auditRepo *repository.AuditRepository,
	jwtSecret string,
	jwtExpiration time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		auditRepo:     auditRepo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest, actorIP string) (*vo.TokenVO, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Create audit log
	s.auditRepo.Create(ctx, &model.AuditLog{
		ActorID:    &user.ID,
		ActorIP:    actorIP,
		Action:     model.ActionLogin,
		ObjectType: model.ObjectTypeUser,
		ObjectID:   &user.ID,
	})

	return &vo.TokenVO{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.jwtExpiration.Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// CreateAPIKey creates a new API key for a user
func (s *AuthService) CreateAPIKey(ctx context.Context, userID uuid.UUID, req dto.CreateAPIKeyRequest, actorIP string) (*vo.APIKeyVO, error) {
	// Generate random API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}
	rawKey := "hive_ak_" + hex.EncodeToString(keyBytes)
	keyHash := repository.HashAPIKey(rawKey)

	// Set default scopes
	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = []string{model.ScopeRead}
	}

	// Calculate expiration
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		exp := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &exp
	}

	apiKey := &model.APIKey{
		UserID:    userID,
		KeyHash:   keyHash,
		Name:      req.Name,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
		IsActive:  true,
	}

	if err := s.userRepo.CreateAPIKey(ctx, apiKey); err != nil {
		return nil, err
	}

	// Create audit log
	s.auditRepo.Create(ctx, &model.AuditLog{
		ActorID:    &userID,
		ActorIP:    actorIP,
		Action:     model.ActionCreate,
		ObjectType: model.ObjectTypeAPIKey,
		ObjectID:   &apiKey.ID,
		Diff: model.JSONMap{
			"name":   apiKey.Name,
			"scopes": apiKey.Scopes,
		},
	})

	return &vo.APIKeyVO{
		ID:        apiKey.ID.String(),
		Name:      apiKey.Name,
		Key:       rawKey, // Only shown once on creation
		Scopes:    apiKey.Scopes,
		ExpiresAt: apiKey.ExpiresAt,
		CreatedAt: apiKey.CreatedAt,
		IsActive:  apiKey.IsActive,
	}, nil
}

// ValidateAPIKey validates an API key and returns the associated user
func (s *AuthService) ValidateAPIKey(ctx context.Context, rawKey string) (*model.User, *model.APIKey, error) {
	keyHash := repository.HashAPIKey(rawKey)

	apiKey, err := s.userRepo.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, nil, err
	}
	if apiKey == nil || !apiKey.IsValid() {
		return nil, nil, ErrInvalidToken
	}

	// Update last used
	s.userRepo.UpdateAPIKeyLastUsed(ctx, apiKey.ID)

	return &apiKey.User, apiKey, nil
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(s.jwtExpiration)

	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "the-hive",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
