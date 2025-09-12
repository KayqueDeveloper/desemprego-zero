package models

import "gorm.io/gorm"

type Partner struct {
	gorm.Model
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	City        string `json:"city" validate:"required"`
	Email       string `json:"email" validate:"required,email" gorm:"unique"`
	Area        Area   `json:"area" validate:"required"`
	Password    string `json:"-"` // Não serializar
}
