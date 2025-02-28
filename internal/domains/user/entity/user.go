package entity

import (
	"zyosa/internal/types"
)

type User struct {
	types.UUID
    Username  string  `json:"username" gorm:"unique;not null"`
    FirstName string  `json:"first_name"`
    LastName  string  `json:"last_name"`
    Email     string  `json:"email" gorm:"unique;not null"`
    Phone     *string `json:"phone,omitempty"`
    Password  string  `json:"password"`
    Image     *string `json:"image,omitempty"`
    IsActive  bool    `json:"is_active"`
    IsDeleted bool    `json:"is_deleted"`
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

func (u *User) TableName() string{
	return "users"
}