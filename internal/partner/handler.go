package partner

import (
	"desemprego-zero/internal/models"
	"desemprego-zero/internal/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// CreatePartner cadastra uma nova empresa parceira
func (h *Handler) CreatePartner(c *gin.Context) {
	var partner models.Partner
	if err := c.ShouldBindJSON(&partner); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	//	if err := validator.ValidaStruct(partner); err != nil {
	//		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
	//		return
	//	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(partner.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}
	partner.Password = string(hashedPassword)

	if err := h.db.Create(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cadastrar empresa"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": partner.ID, "name": partner.Name})
}

// ListPartners lista empresas parceiras (apenas nome e descrição)
func (h *Handler) ListPartners(c *gin.Context) {
	var partners []models.Partner
	if err := h.db.Find(&partners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar empresas"})
		return
	}

	// Retornar apenas nome e descrição
	var result []gin.H
	for _, p := range partners {
		result = append(result, gin.H{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
		})
	}
	c.JSON(http.StatusOK, result)
}
