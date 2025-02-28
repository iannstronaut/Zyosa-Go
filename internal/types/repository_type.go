package types

import (
	"gorm.io/gorm"
)

type Repository[T any] struct{
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

func (r *Repository[T]) FindByID(id T, entity *T) error {
	return r.DB.Where("id = ?", id).First(entity).Error
}

func (r *Repository[T]) FindAll(limit int, offset int, entities *[]T) error {
	return r.DB.Limit(limit).Offset(offset).Find(entities).Error
}

func (r *Repository[T]) CountByID(id string) (int64, error) {
	var count int64
	err := r.DB.Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count, err
}




