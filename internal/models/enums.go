package models

// Status do candidato

// CandidateStatus representa a situação atual do candidato
// Exemplo: empregado, desempregado, em busca

//go:generate stringer -type=CandidateStatus

type CandidateStatus string

const (
	StatusEmBusca      CandidateStatus = "em_busca"
	StatusEmpregado    CandidateStatus = "empregado"
	StatusDesempregado CandidateStatus = "desempregado"
)

// Área de atuação

type Area string

const (
	AreaTI             Area = "TI"
	AreaSaude          Area = "Saude"
	AreaAdministrativa Area = "Administrativa"
	// ... outras áreas
)

// Tipo de vaga

type JobType string

const (
	JobFullTime   JobType = "Full-time"
	JobPartTime   JobType = "Part-time"
	JobContract   JobType = "Contract"
	JobInternship JobType = "Internship"
)
