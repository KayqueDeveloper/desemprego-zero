package models

import (
	"gorm.io/gorm"
)

type Candidate struct {
	gorm.Model
	Name       string `json:"name" validate:"required,min=3,max=100"`
	Email      string `json:"email" validate:"required,email" gorm:"uniqueIndex:idx_candidate_job"`
	Phone      string `json:"phone" validate:"required"`
	Resume     string `json:"resume" validate:"required,min=10"`
	Experience string `json:"experience" validate:"required,min=10"`
	Education  string `json:"education" validate:"required,min=10"`
	Jobs       []Job  `json:"jobs" gorm:"many2many:job_candidates;uniqueIndex:idx_candidate_job"`
}
