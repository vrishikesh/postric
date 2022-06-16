package models

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const dbTimeout = time.Second * 3

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

type Models struct {
	User
}

type User struct {
	gorm.Model
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
}

func (u *User) Insert() (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return 0, err
	}

	newUser := *u
	newUser.Password = string(hashedPassword)
	result := db.WithContext(ctx).Create(&newUser)
	if result.Error != nil {
		return 0, result.Error
	}

	return newUser.ID, nil
}

func (u *User) GetAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var users []User
	result := db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
