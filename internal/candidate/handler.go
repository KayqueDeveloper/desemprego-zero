package candidate

import (
	"desemprego-zero/internal/models"
	"desemprego-zero/internal/validator"
	"net/http"
	"strings"

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

// CreateCandidate cadastra um novo candidato com upload de currículo
func (h *Handler) CreateCandidate(c *gin.Context) {
	var candidate models.Candidate

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário"})
		return
	}

	candidate.Name = c.PostForm("name")
	candidate.Email = c.PostForm("email")
	candidate.Password = c.PostForm("password")
	candidate.Phone = c.PostForm("phone")
	candidate.City = c.PostForm("city")
	candidate.Status = models.CandidateStatus(c.PostForm("status"))
	candidate.Area = models.Area(c.PostForm("area"))
	candidate.Experience = c.PostForm("experience")

	// Validação
	if err := validator.ValidaStruct(candidate); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	if !validator.ValidaEmail(candidate.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		return
	}

	//	if !validator.ValidaTelefone(candidate.Phone) {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "Telefone inválido. Use o formato (99) 99999-9999"})
	//		return
	//	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(candidate.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}
	candidate.Password = string(hashedPassword)

	// // Upload do currículo
	// file, header, err := c.Request.FormFile("resume")
	// if err == nil {
	// 	defer file.Close()
	// 	filename := filepath.Base(header.Filename)
	// 	resumePath := filepath.Join("uploads", "resumes", filename)
	// 	os.MkdirAll(filepath.Dir(resumePath), 0755)
	// 	out, err := os.Create(resumePath)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar currículo"})
	// 		return
	// 	}
	// 	defer out.Close()
	// 	_, err = out.ReadFrom(file)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gravar currículo"})
	// 		return
	// 	}
	// 	candidate.Resume = resumePath
	// }

	if err := h.db.Create(&candidate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cadastrar candidato"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": candidate.ID, "name": candidate.Name})
}

// ListCandidates lista candidatos com filtros por área, status e nome
func (h *Handler) ListCandidates(c *gin.Context) {
	var candidates []models.Candidate
	area := c.Query("area")
	status := c.Query("status")
	name := c.Query("name")

	dbQuery := h.db.Model(&models.Candidate{})
	if area != "" {
		dbQuery = dbQuery.Where("area = ?", area)
	}
	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}
	if name != "" {
		dbQuery = dbQuery.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}
	dbQuery = dbQuery.Preload("Jobs").Find(&candidates)

	c.JSON(http.StatusOK, candidates)
}
