package models

import "gorm.io/gorm"

type Feedback struct {
	gorm.Model
	CandidateID uint   `json:"candidate_id" validate:"required"`
	Message     string `json:"message" validate:"required"`
	PartnerID   uint   `json:"partner_id" validate:"required"`
}
