package types

import (
	"errors"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *Repository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *Repository[T]) Delete(entity *T) error {
	return r.DB.Delete(entity).Error
}

func (r *Repository[T]) FindByID(id string, entity *T) (*T, error) {
	err := r.DB.Where("id = ?", id).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return entity, nil
}

func (r *Repository[T]) FindAll(limit int, offset int, entities *[]T) (*[]T, error) {
	err := r.DB.Limit(limit).Offset(offset).Find(entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return entities, nil
}

func (r *Repository[T]) CountByID(id string) (*int64, error) {
	var count int64
	err := r.DB.Model(new(T)).Where("id = ?", id).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &count, err
}
