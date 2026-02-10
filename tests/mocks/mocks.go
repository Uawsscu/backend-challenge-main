package mocks

import (
	"context"
	"time"

	"github.com/backend-challenge/user-api/internal/domain"
)

type MockUserRepository struct {
	CreateFunc      func(ctx context.Context, user *domain.User) error
	FindByIDFunc    func(ctx context.Context, id string) (*domain.User, error)
	FindByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	FindAllFunc     func(ctx context.Context) ([]*domain.User, error)
	UpdateFunc      func(ctx context.Context, user *domain.User) error
	DeleteFunc      func(ctx context.Context, id string) error
	CountFunc       func(ctx context.Context) (int64, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return []*domain.User{}, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

type MockSessionManager struct {
	StoreSessionFunc  func(ctx context.Context, key string, data interface{}, ttl time.Duration) error
	GetSessionFunc    func(ctx context.Context, key string) (string, error)
	DeleteSessionFunc func(ctx context.Context, key string) error
}

func (m *MockSessionManager) StoreSession(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	if m.StoreSessionFunc != nil {
		return m.StoreSessionFunc(ctx, key, data, ttl)
	}
	return nil
}

func (m *MockSessionManager) GetSession(ctx context.Context, key string) (string, error) {
	if m.GetSessionFunc != nil {
		return m.GetSessionFunc(ctx, key)
	}
	return "", nil
}

func (m *MockSessionManager) DeleteSession(ctx context.Context, key string) error {
	if m.DeleteSessionFunc != nil {
		return m.DeleteSessionFunc(ctx, key)
	}
	return nil
}

type MockTokenService struct {
	GenerateTokenFunc func(userID, email string, duration time.Duration) (string, error)
	ValidateTokenFunc func(tokenString string) (*domain.TokenClaims, error)
}

func (m *MockTokenService) GenerateToken(userID, email string, duration time.Duration) (string, error) {
	if m.GenerateTokenFunc != nil {
		return m.GenerateTokenFunc(userID, email, duration)
	}
	return "mock-token", nil
}

func (m *MockTokenService) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(tokenString)
	}
	return &domain.TokenClaims{
		Subject: "test-token-id",
	}, nil
}

type MockLotteryRepository struct {
	SearchAndReserveFunc func(ctx context.Context, pattern string, userID string, limit int) ([]domain.LotteryTicket, error)
	UpsertManyFunc       func(ctx context.Context, tickets []domain.LotteryTicket) error
	MarkAsSoldFunc       func(ctx context.Context, ticketID string, userID string) error
	CountFunc            func(ctx context.Context) (int64, error)
	SeedTicketsFunc      func(ctx context.Context, total int) error
}

func (m *MockLotteryRepository) SearchAndReserve(ctx context.Context, pattern string, userID string, limit int) ([]domain.LotteryTicket, error) {
	if m.SearchAndReserveFunc != nil {
		return m.SearchAndReserveFunc(ctx, pattern, userID, limit)
	}
	return []domain.LotteryTicket{}, nil
}

func (m *MockLotteryRepository) UpsertMany(ctx context.Context, tickets []domain.LotteryTicket) error {
	if m.UpsertManyFunc != nil {
		return m.UpsertManyFunc(ctx, tickets)
	}
	return nil
}

func (m *MockLotteryRepository) MarkAsSold(ctx context.Context, ticketID string, userID string) error {
	if m.MarkAsSoldFunc != nil {
		return m.MarkAsSoldFunc(ctx, ticketID, userID)
	}
	return nil
}

func (m *MockLotteryRepository) Count(ctx context.Context) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

func (m *MockLotteryRepository) SeedTickets(ctx context.Context, total int) error {
	if m.SeedTicketsFunc != nil {
		return m.SeedTicketsFunc(ctx, total)
	}
	return nil
}
