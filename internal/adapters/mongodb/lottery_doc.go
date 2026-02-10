package mongodb

import (
	"time"

	"github.com/backend-challenge/user-api/internal/domain"
)

type lotteryDoc struct {
	ID            string               `bson:"_id,omitempty"`
	Number        string               `bson:"number"`
	Status        domain.LotteryStatus `bson:"status"`
	ReservedUntil *time.Time           `bson:"reserved_until,omitempty"`
	ReservedBy    string               `bson:"reserved_by,omitempty"`
	CreatedAt     time.Time            `bson:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at"`
}

func fromLotteryDomain(l *domain.LotteryTicket) *lotteryDoc {
	if l == nil {
		return nil
	}
	return &lotteryDoc{
		ID:            l.ID,
		Number:        l.Number,
		Status:        l.Status,
		ReservedUntil: l.ReservedUntil,
		ReservedBy:    l.ReservedBy,
		CreatedAt:     l.CreatedAt,
		UpdatedAt:     l.UpdatedAt,
	}
}

func (d *lotteryDoc) toLotteryDomain() *domain.LotteryTicket {
	if d == nil {
		return nil
	}
	return &domain.LotteryTicket{
		ID:            d.ID,
		Number:        d.Number,
		Status:        d.Status,
		ReservedUntil: d.ReservedUntil,
		ReservedBy:    d.ReservedBy,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
