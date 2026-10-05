package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/seu-usuario/taskflow-backend/internal/database"
	"github.com/seu-usuario/taskflow-backend/internal/logger"
	"github.com/seu-usuario/taskflow-backend/internal/routes"
)

func main() {
	logger.Init()
	if err := run(); err != nil {
		logger.Errorf("falha ao iniciar a API: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Carrega o .env local; variáveis já exportadas pelo ambiente têm precedência.
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Warnf("não foi possível carregar o arquivo .env: %v", err)
	}

	// Cria e valida a conexão com o banco de dados antes de iniciar o servidor.
	db, err := database.Connect(database.Config{
		Host:     envOrDefault("DB_HOST", "localhost"),
		Port:     envOrDefault("DB_PORT", "3306"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	})
	if err != nil {
		return fmt.Errorf("conectar ao banco de dados: %w", err)
	}
	defer db.Close()

	// Compõe as rotas e suas dependências, depois inicia o servidor HTTP.
	router := routes.SetupRouter(db)
	address := ":" + envOrDefault("PORT", "8080")
	logger.Infof("servidor Gin iniciado em %s", address)
	if err := router.Run(address); err != nil {
		return fmt.Errorf("iniciar servidor HTTP: %w", err)
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
