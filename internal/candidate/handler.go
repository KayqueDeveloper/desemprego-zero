package candidate

import (
	"net/http"

	"desemprego-zero/internal/models"
	"desemprego-zero/internal/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// CreateCandidate permite que um candidato se inscreva em uma vaga
func (h *Handler) CreateCandidate(c *gin.Context) {
	var candidate models.Candidate
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	// Valida a estrutura
	if err := validator.ValidaStruct(candidate); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	// Valida o email
	if !validator.ValidaEmail(candidate.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		return
	}

	// Valida o telefone
	if !validator.ValidaTelefone(candidate.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telefone inválido. Use o formato (99) 99999-9999"})
		return
	}

	// Verifica se a vaga existe
	var job models.Job
	if err := h.db.First(&job, candidate.Jobs[0].ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaga não encontrada"})
		return
	}

	// Verifica se o candidato já se inscreveu nesta vaga
	var existingCandidate models.Candidate
	result := h.db.Where("email = ? AND jobs.id = ?", candidate.Email, job.ID).
		Joins("JOIN job_candidates ON job_candidates.candidate_id = candidates.id").
		Joins("JOIN jobs ON jobs.id = job_candidates.job_id").
		First(&existingCandidate)

	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Você já se inscreveu nesta vaga"})
		return
	}

	// Cria o candidato e associa à vaga
	tx := h.db.Begin()
	if err := tx.Create(&candidate).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar candidato"})
		return
	}

	if err := tx.Model(&candidate).Association("Jobs").Append(&job); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao associar candidato à vaga"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, candidate)
}

// ListCandidates lista todos os candidatos (rota protegida)
func (h *Handler) ListCandidates(c *gin.Context) {
	var candidates []models.Candidate
	result := h.db.Preload("Jobs").Find(&candidates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar candidatos"})
		return
	}

	c.JSON(http.StatusOK, candidates)
}
