package feedback

import (
	"desemprego-zero/internal/models"
	"desemprego-zero/internal/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// CreateFeedback cadastra um novo feedback
func (h *Handler) CreateFeedback(c *gin.Context) {
	var feedback models.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	if err := validator.ValidaStruct(feedback); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	// Verifica se o candidato existe
	var candidate models.Candidate
	if err := h.db.First(&candidate, feedback.CandidateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Candidato não encontrado"})
		return
	}

	// Verifica se a empresa existe
	var partner models.Partner
	if err := h.db.First(&partner, feedback.PartnerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa não encontrada"})
		return
	}

	if err := h.db.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cadastrar feedback"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": feedback.ID})
}
