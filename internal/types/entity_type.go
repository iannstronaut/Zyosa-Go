package types

import (
	"github.com/google/uuid"

	"time"
)

type UUID struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v7();index" json:"id"`
}

type IdIncrement struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}