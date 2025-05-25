package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"` // O "-" faz com que o campo não seja serializado em JSON
	Email    string `json:"email" gorm:"unique"`
}
