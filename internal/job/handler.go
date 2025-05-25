package job

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

// ListJobs retorna todas as vagas ativas
func (h *Handler) ListJobs(c *gin.Context) {
	var jobs []models.Job
	result := h.db.Where("active = ?", true).Find(&jobs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar vagas"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// GetJob retorna os detalhes de uma vaga específica
func (h *Handler) GetJob(c *gin.Context) {
	id := c.Param("id")
	var job models.Job

	result := h.db.First(&job, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaga não encontrada"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// CreateJob cria uma nova vaga (rota protegida)
func (h *Handler) CreateJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	// Valida a estrutura
	if err := validator.ValidaStruct(job); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	result := h.db.Create(&job)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar vaga"})
		return
	}

	c.JSON(http.StatusCreated, job)
}

// UpdateJob atualiza uma vaga existente (rota protegida)
func (h *Handler) UpdateJob(c *gin.Context) {
	id := c.Param("id")
	var job models.Job

	if err := h.db.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaga não encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	// Valida a estrutura
	if err := validator.ValidaStruct(job); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	result := h.db.Save(&job)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar vaga"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// DeleteJob exclui uma vaga (rota protegida)
func (h *Handler) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	var job models.Job

	if err := h.db.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaga não encontrada"})
		return
	}

	result := h.db.Delete(&job)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir vaga"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vaga excluída com sucesso"})
}
