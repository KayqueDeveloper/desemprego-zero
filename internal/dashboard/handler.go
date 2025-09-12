package dashboard

import (
	"desemprego-zero/internal/models"
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

// CountDashboard retorna os contadores principais do sistema
func (h *Handler) CountDashboard(c *gin.Context) {
	var countEmpregados int64
	h.db.Model(&models.Candidate{}).Where("status = ?", models.StatusEmpregado).Count(&countEmpregados)

	var countEmpresas int64
	h.db.Model(&models.Partner{}).Count(&countEmpresas)

	var countFeedbacks int64
	h.db.Model(&models.Feedback{}).Count(&countFeedbacks)

	taxaSatisfacao := 0.0
	if countEmpregados > 0 {
		taxaSatisfacao = (float64(countFeedbacks) / float64(countEmpregados)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"pessoas_empregadas": countEmpregados,
		"empresas_parceiras": countEmpresas,
		"taxa_satisfacao":    taxaSatisfacao,
	})
}
