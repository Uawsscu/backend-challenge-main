package ports

import (
	"context"

	"github.com/backend-challenge/user-api/internal/domain"
)

type LotteryRepository interface {
	SearchAndReserve(ctx context.Context, pattern string, userID string, limit int) ([]domain.LotteryTicket, error)
	UpsertMany(ctx context.Context, tickets []domain.LotteryTicket) error
	MarkAsSold(ctx context.Context, ticketID string, userID string) error
	Count(ctx context.Context) (int64, error)
	SeedTickets(ctx context.Context, total int) error
}
