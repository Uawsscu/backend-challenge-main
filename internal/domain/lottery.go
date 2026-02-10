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
	ID            string        `json:"id"`
	Number        string        `json:"number"` // 6-digit number
	Status        LotteryStatus `json:"status"`
	ReservedUntil *time.Time    `json:"reserved_until,omitempty"`
	ReservedBy    string        `json:"reserved_by,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
