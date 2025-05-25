package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title        string      `json:"title" validate:"required,min=3,max=100"`
	Description  string      `json:"description" validate:"required,min=10"`
	Company      string      `json:"company" validate:"required,min=2,max=100"`
	Location     string      `json:"location" validate:"required,min=3,max=100"`
	Salary       string      `json:"salary" validate:"required"`
	Type         string      `json:"type" validate:"required,oneof=Full-time Part-time Contract Internship"`
	Requirements string      `json:"requirements" validate:"required,min=10"`
	Deadline     time.Time   `json:"deadline" validate:"required,gtfield=CreatedAt"`
	Active       bool        `json:"active" gorm:"default:true"`
	Candidates   []Candidate `json:"candidates" gorm:"many2many:job_candidates;"`
}
