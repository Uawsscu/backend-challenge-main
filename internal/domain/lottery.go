package domain

import (
	"time"
)

type LotteryStatus string

const (
	LotteryStatusAvailable LotteryStatus = "available"
	LotteryStatusReserved  LotteryStatus = "reserved"
	LotteryStatusSold      LotteryStatus = "sold"
)

type LotteryTicket struct {
	ID            string
	Number        string
	Status        LotteryStatus
	ReservedUntil *time.Time
	ReservedBy    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
