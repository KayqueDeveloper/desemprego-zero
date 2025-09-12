package models

import "gorm.io/gorm"

type Candidate struct {
	gorm.Model
	Name       string          `json:"name" validate:"required,min=3,max=100"`
	Email      string          `json:"email" validate:"required,email" gorm:"unique"`
	Password   string          `json:"-"` // Para autenticação
	Phone      string          `json:"phone"`
	City       string          `json:"city"`
	Status     CandidateStatus `json:"status"`
	Area       Area            `json:"area"`
	Experience string          `json:"experience"`
	Resume     string          `json:"resume"` // Caminho do arquivo
	Jobs       []Job           `json:"jobs" gorm:"many2many:job_candidates;"`
}
