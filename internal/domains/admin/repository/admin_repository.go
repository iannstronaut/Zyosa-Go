package repository

import (
	"zyosa/internal/domains/admin/entity"
	"zyosa/internal/types"

	"gorm.io/gorm"
)

type AdminRepository struct {
	types.Repository[entity.Admin]
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{Repository: types.Repository[entity.Admin]{DB: db}}
}

func (r *AdminRepository) FindByUsername(admin *entity.Admin) (*entity.Admin, error) {
	var exists entity.Admin
	if err := r.DB.Where("username = ?", admin.Username).First(&exists).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &exists, nil
}

func (r *AdminRepository) FindByEmail(admin *entity.Admin) (*entity.Admin, error) {
	var exists entity.Admin
	if err := r.DB.Where("email = ?", admin.Email).First(&exists).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &exists, nil
}
