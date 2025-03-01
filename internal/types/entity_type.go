package types

import (
	"github.com/google/uuid"

	"time"
)

type UUID struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
}

type IdIncrement struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at"`
}
