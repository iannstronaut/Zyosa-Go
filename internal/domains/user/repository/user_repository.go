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