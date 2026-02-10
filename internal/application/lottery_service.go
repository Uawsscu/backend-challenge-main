package application

import (
	"context"
	"fmt"

	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/internal/ports"
)

type LotteryService struct {
	repo ports.LotteryRepository
}

func NewLotteryService(repo ports.LotteryRepository) *LotteryService {
	return &LotteryService{
		repo: repo,
	}
}

func (s *LotteryService) SearchLottery(ctx context.Context, pattern string, userID string) ([]domain.LotteryTicket, error) {
	if len(pattern) != 6 {
		return nil, fmt.Errorf("%w: pattern must be 6 characters", domain.ErrRequestInvalid)
	}

	for _, char := range pattern {
		if (char < '0' || char > '9') && char != '*' {
			return nil, fmt.Errorf("%w: pattern can only contain digits and *", domain.ErrRequestInvalid)
		}
	}

	return s.repo.SearchAndReserve(ctx, pattern, userID, 10) // limit 10
}

func (s *LotteryService) GetLotteryCount(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
