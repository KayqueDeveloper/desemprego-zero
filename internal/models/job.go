package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title        string      `json:"title" validate:"required"`
	PartnerID    uint        `json:"partner_id" validate:"required"`
	Partner      Partner     `json:"partner"`
	Location     string      `json:"location" validate:"required"`
	Salary       float64     `json:"salary" validate:"required"`
	Type         JobType     `json:"type" validate:"required"`
	Deadline     string      `json:"deadline" validate:"required"`
	Description  string      `json:"description" validate:"required"`
	Requirements string      `json:"requirements" validate:"required"`
	Candidates   []Candidate `json:"candidates" gorm:"many2many:job_candidates;"`
}
