package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/backend-challenge/user-api/internal/application"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/tests/mocks"
)

func TestLotteryService_SearchLottery(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		userID      string
		mockSetup   func(*mocks.MockLotteryRepository)
		expectError bool
		errorType   error
	}{
		{
			name:    "successful search",
			pattern: "123***",
			userID:  "user-123",
			mockSetup: func(repo *mocks.MockLotteryRepository) {
				repo.SearchAndReserveFunc = func(ctx context.Context, pattern string, userID string, limit int) ([]domain.LotteryTicket, error) {
					return []domain.LotteryTicket{
						{Number: "123456"},
						{Number: "123000"},
					}, nil
				}
			},
			expectError: false,
		},
		{
			name:        "invalid pattern length",
			pattern:     "123",
			userID:      "user-123",
			mockSetup:   func(repo *mocks.MockLotteryRepository) {},
			expectError: true,
		},
		{
			name:        "invalid characters in pattern",
			pattern:     "123-bc",
			userID:      "user-123",
			mockSetup:   func(repo *mocks.MockLotteryRepository) {},
			expectError: true,
		},
		{
			name:    "repository error",
			pattern: "******",
			userID:  "user-123",
			mockSetup: func(repo *mocks.MockLotteryRepository) {
				repo.SearchAndReserveFunc = func(ctx context.Context, pattern string, userID string, limit int) ([]domain.LotteryTicket, error) {
					return nil, errors.New("db error")
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockLotteryRepository{}
			tt.mockSetup(mockRepo)

			service := application.NewLotteryService(mockRepo)
			results, err := service.SearchLottery(context.Background(), tt.pattern, tt.userID)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(results) == 0 {
					t.Errorf("expected results but got none")
				}
			}
		})
	}
}

func TestLotteryService_GetLotteryCount(t *testing.T) {
	mockRepo := &mocks.MockLotteryRepository{}
	service := application.NewLotteryService(mockRepo)

	t.Run("successful count", func(t *testing.T) {
		mockRepo.CountFunc = func(ctx context.Context) (int64, error) {
			return 100, nil
		}

		count, err := service.GetLotteryCount(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count != 100 {
			t.Errorf("expected count 100 but got %d", count)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.CountFunc = func(ctx context.Context) (int64, error) {
			return 0, errors.New("db error")
		}

		_, err := service.GetLotteryCount(context.Background())
		if err == nil {
			t.Errorf("expected error but got none")
		}
	})
}
