package entity

import (
	"zyosa/internal/types"
)

type User struct {
	types.UUID
	Username  string  `json:"username" gorm:"type:varchar(32);unique;not null"`
	FirstName string  `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName  string  `json:"last_name" gorm:"type:varchar(100);not null"`
	Email     string  `json:"email" gorm:"type:varchar(255);unique;not null"`
	Phone     *string `json:"phone,omitempty" gorm:"type:varchar(20)"`
	Password  string  `json:"password" gorm:"type:text;not null"`
	Image     *string `json:"image,omitempty" gorm:"type:text"`
	IsActive  bool    `json:"is_active" gorm:"default:true"`
	IsDeleted bool    `json:"is_deleted" gorm:"default:false"`
	types.Timestamps
}

func NewUser(usename string, firstName string, lastName string, email string, password string, phone *string, image *string) User {
	return User{
		Username:  usename,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Password:  password,
		Image:     image,
		IsActive:  true,
		IsDeleted: false,
	}
}

func (u *User) TableName() string {
	return "users"
}
