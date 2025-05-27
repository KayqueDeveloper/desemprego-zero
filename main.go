package main

import (
	"fmt"
	"log"
	"os"

	"desemprego-zero/internal/auth"
	"desemprego-zero/internal/candidate"
	"desemprego-zero/internal/job"
	"desemprego-zero/internal/middleware"
	"desemprego-zero/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func validateEnvVars() error {
	requiredVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"JWT_SECRET",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("variável de ambiente %s não está definida", v)
		}
	}

	return nil
}

func main() {
	// Tenta carregar o arquivo .env, mas não falha se não existir
	_ = godotenv.Load()

	// Valida variáveis de ambiente
	if err := validateEnvVars(); err != nil {
		log.Fatalf("Erro nas variáveis de ambiente: %v", err)
	}

	// Configuração do banco de dados
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Migra as tabelas
	err = db.AutoMigrate(
		&models.Job{},
		&models.Candidate{},
		&models.Admin{},
	)
	if err != nil {
		log.Fatalf("Erro ao migrar tabelas: %v", err)
	}

	// Inicializa o router
	r := gin.Default()

	// Middleware de tratamento de erros
	r.Use(middleware.ErrorHandler())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"db":     "connected",
		})
	})

	// Inicializa os handlers
	jobHandler := job.NewHandler(db)
	candidateHandler := candidate.NewHandler(db)
	authHandler := auth.NewHandler(db)

	// Rotas públicas
	public := r.Group("/")
	{
		// Rotas de vagas
		public.GET("/jobs", jobHandler.ListJobs)
		public.GET("/jobs/:id", jobHandler.GetJob)
		public.POST("/candidates", candidateHandler.CreateCandidate)
	}

	// Rotas protegidas (admin)
	admin := r.Group("/admin")
	{
		admin.POST("/login", authHandler.Login)
	}

	// Rotas que requerem autenticação
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Rotas de vagas protegidas
		protected.POST("/jobs", jobHandler.CreateJob)
		protected.PUT("/jobs/:id", jobHandler.UpdateJob)
		protected.DELETE("/jobs/:id", jobHandler.DeleteJob)

		// Rotas de candidatos protegidas
		protected.GET("/candidates", candidateHandler.ListCandidates)
	}

	// Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
