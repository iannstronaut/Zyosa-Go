package repository

import (
	"zyosa/internal/domains/user/entity"
	"zyosa/internal/types"

	"gorm.io/gorm"
)

type UserRepository struct {
	types.Repository[entity.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Repository: types.Repository[entity.User]{DB: db}}
}

func (r *UserRepository) FindByUsername(user *entity.User) (*entity.User, error) {
	var exists entity.User
	if err := r.DB.Where("username = ?", user.Username).First(&exists).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err 
		}
		return nil, err 
	}
	return &exists, nil 
}

func (r *UserRepository) FindByEmail(user *entity.User) (*entity.User, error) {
	var exists entity.User
	if err := r.DB.Where("email = ?", user.Email).First(&exists).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil 
		}
		return nil, err
	}
	return &exists, nil 
}

