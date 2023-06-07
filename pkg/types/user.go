package types

import "time"

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName  string    `json:"username" binding:"required,min=3,max=255" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" binding:"required,email" gorm:"type:varchar(255);unique; not null"`
	Password  string    `json:"password" binding:"required,min=6,max=16" gorm:"type:varchar(60); not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
