package models

import (
	"encoding/json"
	"time"
)

// fmt.Stringer implementation of a Block
func (i Block) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

// fmt.Stringer implementation of a Slot
func (i Slot) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

// fmt.Stringer implementation of a Epoch
func (i Epoch) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

// represents a block
type Block struct {
	BlockNumber      uint64    `json:"blockNumber" db:"block_number"`
	BlockRoot        string    `json:"blockRoot" db:"block_root"`
	StateRoot        string    `json:"stateRoot" db:"state_root"`
	SlotNumber       uint64    `json:"slotNumber" db:"slot_number"`
	GasLimit         uint64    `json:"gasLimit" db:"gas_limit"`
	GasUsed          uint64    `json:"gasUsed" db:"gas_used"`
	NoOfTransactions int       `json:"noOfTransactions" db:"no_of_transactions"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// represents a slot
type Slot struct {
	SlotNumber  uint64    `json:"slotNumber" db:"slot_number"`
	StartTime   time.Time `json:"startTime" db:"start_time"`
	EndTime     time.Time `json:"endTime" db:"end_time"`
	EpochNumber uint64    `json:"epochNumber" db:"epoch_number"`
	Block       Block     `json:"block"`
}

// represents an epoch
type Epoch struct {
	EpochNumber uint64    `json:"epochNumber" db:"epoch_number"`
	StartTime   time.Time `json:"startTime" db:"start_time"`
	EndTime     time.Time `json:"endTime" db:"end_time"`
	Slots       []Slot    `json:"slots"`
}
