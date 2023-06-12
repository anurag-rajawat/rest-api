package types

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName  string  `json:"username" gorm:"type:varchar(255);not null"`
	Email     string  `json:"email" gorm:"type:varchar(255);unique; not null"`
	Password  string  `json:"password" gorm:"type:varchar(60); not null"`
	CreatedAt *string `json:"created_at,omitempty" gorm:"type:varchar(10)"`
}

func (u *User) Create(db *gorm.DB) (*User, error) {
	hashedPasswd, err := hashPassword(u.Password)
	if err != nil {
		return &User{}, err
	}
	u.Password = hashedPasswd

	err = db.Model(&User{}).Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAll(db *gorm.DB) (*[]User, error) {
	var users []User
	err := db.Model(&User{}).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, nil
}

func (u *User) FindById(db *gorm.DB, id uint64) (*User, error) {
	var user User
	err := db.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("user not found")
	}
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

func (u *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.Model(&User{}).Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
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
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetJwtToken() (string, error) {
	claims := jwt.MapClaims{
		"id":         u.ID,
		"email":      u.Email,
		"expires_at": time.Now().Add(24 * time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Error(err)
		return "", err
	}
	return token, nil
}

func (u *User) CheckPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func hashPassword(password string) (string, error) {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswd), nil
}
