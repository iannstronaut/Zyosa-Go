package entity

import "zyosa/internal/types"

type Session struct {
	types.IdIncrement
	UserId    string `json:"user_id" gorm:"type:char(36);not null;index:idx_user_token"`
	Token     string `json:"token" gorm:"type:char(36);not null;index:idx_user_token"`
	UserAgent string `json:"user_agent" gorm:"type:text;not null"`
	ExpiredAt string `json:"expired_at" gorm:"type:datetime;not null"`
	UsedAt    string `json:"used_at" gorm:"type:datetime;default:null"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
}

func NewSession(userId string, token string, userAgent string, expiredAt string) Session {
	return Session{
		UserId:    userId,
		Token:     token,
		UserAgent: userAgent,
		ExpiredAt: expiredAt,
		IsActive:  true,
	}
}

func (s *Session) TableName() string {
	return "session"
}