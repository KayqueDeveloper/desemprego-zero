package main

import (
	"log"
	"os"

	"desemprego-zero/internal/auth"
	"desemprego-zero/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}

	// Configuração do banco de dados
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Migra as tabelas
	db.AutoMigrate(&models.Admin{})

	// Cria o handler de autenticação
	authHandler := auth.NewHandler(db)

	// Cria o administrador inicial
	err = authHandler.CreateAdmin(
		"admin",
		"admin123",
		"admin@igreja.com",
	)
	if err != nil {
		log.Fatal("Erro ao criar administrador:", err)
	}

	log.Println("Administrador criado com sucesso!")
}
