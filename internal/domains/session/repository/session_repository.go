package repository

import (
	"zyosa/internal/domains/session/entity"
	"zyosa/internal/types"

	"gorm.io/gorm"
)

type SessionRepository struct {
	types.Repository[entity.Session]
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	db.AutoMigrate(&entity.Session{})
	return &SessionRepository{Repository: types.Repository[entity.Session]{DB: db}}
}

func (s *SessionRepository) CreateSession(session *entity.Session) error {
	if err := s.DB.Create(session).Error; err != nil {
		return err
	}
	return nil
}

func (s *SessionRepository) FindByToken(token string) (*entity.Session, error) {
	var exist entity.Session
	if err := s.DB.Where("token = ?", token).First(&exist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &exist, nil
}

func (s *SessionRepository) DeactiveSession(session *entity.Session) error {
	if err := s.DB.Model(session).Update("is_active", false).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}
	return nil
}