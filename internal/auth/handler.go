package auth

import (
	"net/http"
	"os"
	"time"

	"desemprego-zero/internal/models"
	"desemprego-zero/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" validate:"required,min=6"`
}

// Login autentica admin, candidato ou empresa parceira e retorna um token JWT
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	if err := validator.ValidaStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErroValidacaoParaResponse(err))
		return
	}

	var (
		userID   uint
		role     string
		userData gin.H
	)

	// Tenta admin
	var admin models.Admin
	if err := h.db.Where("username = ?", req.Username).Or("email = ?", req.Email).First(&admin).Error; err == nil {
		if bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)) == nil {
			userID = admin.ID
			role = "admin"
			userData = gin.H{"id": admin.ID, "username": admin.Username, "email": admin.Email}
		}
	}

	// Tenta candidate
	if userID == 0 {
		var candidate models.Candidate
		if err := h.db.Where("email = ?", req.Email).First(&candidate).Error; err == nil {
			if bcrypt.CompareHashAndPassword([]byte(candidate.Password), []byte(req.Password)) == nil {
				userID = candidate.ID
				role = "candidate"
				userData = gin.H{"id": candidate.ID, "name": candidate.Name, "email": candidate.Email}
			}
		}
	}

	// Tenta partner
	if userID == 0 {
		var partner models.Partner
		if err := h.db.Where("email = ?", req.Email).First(&partner).Error; err == nil {
			if bcrypt.CompareHashAndPassword([]byte(partner.Password), []byte(req.Password)) == nil {
				userID = partner.ID
				role = "partner"
				userData = gin.H{"id": partner.ID, "name": partner.Name, "email": partner.Email}
			}
		}
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Gera o token JWT com expiração de 24 horas
	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "desemprego-zero",
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_at": expirationTime,
		"role":       role,
		"user":       userData,
	})
}

// CreateAdmin cria um novo administrador (função auxiliar para setup inicial)
func (h *Handler) CreateAdmin(username, password, email string) error {
	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.Admin{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	return h.db.Create(&admin).Error
}
