package types

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName  string    `json:"username" binding:"required,min=3,max=255" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" binding:"required,email" gorm:"type:varchar(255);unique; not null"`
	Password  string    `json:"password" binding:"required,min=6,max=16" gorm:"type:varchar(60); not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (u *User) FindAll(db *gorm.DB) (*[]User, error) {
	var users []User
	err := db.Model(&User{}).Find(&users).Error
	if err != nil {
		return &[]User{}, nil
	}
	return &users, nil
}

func (u *User) FindById(db *gorm.DB, id uint64) (*User, error) {
	var user User
	err := db.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return &User{}, errors.New("user not found")
	}
	return &user, nil
}

func (u *User) Update(db *gorm.DB, id uint64) (*User, error) {
	var user User
	_, err := user.FindById(db, id)
	if err != nil {
		return &User{}, err
	}

	hashedPasswd, err := hashPassword(u.Password)
	if err != nil {
		return &User{}, err
	}

	u.Password = hashedPasswd
	u.ID = id
	err = db.Where("id = ?", id).First(&user).Updates(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) Delete(db *gorm.DB, id uint64) error {
	_, err := u.FindById(db, id)
	if err != nil {
		return err
	}

	err = db.Model(&User{}).Where("id = ?", id).Delete(&u).Error
	if err == gorm.ErrRecordNotFound {
		return errors.New("user not found")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswd), nil
}
