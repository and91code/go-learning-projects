package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/seu-usuario/taskflow-backend/internal/database"
	"github.com/seu-usuario/taskflow-backend/internal/logger"
	"github.com/seu-usuario/taskflow-backend/internal/routes"
)

// Config centraliza as variáveis de ambiente necessárias
type Config struct {
	Port       string
	JWTSecret  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBDSN      string
	RedisURL   string
}

func loadConfig() Config {
	return Config{
		Port:       envOrDefault("PORT", "8080"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		DBHost:     envOrDefault("DB_HOST", "localhost"),
		DBPort:     envOrDefault("DB_PORT", "3306"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBDSN:      os.Getenv("DB_DSN"),
		RedisURL:   envOrDefault("REDIS_URL", "redis://localhost:6379"),
	}
}

func main() {
	logger.Init()

	if err := run(); err != nil {
		logger.Errorf("falha ao iniciar a API: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Carrega o .env se existir
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Warnf("não foi possível carregar o arquivo .env: %v", err)
	}

	cfg := loadConfig()

	// 2. Conexão com o banco de dados
	db, err := database.Connect(database.Config{
		DSN:      cfg.DBDSN,
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
	})
	if err != nil {
		return fmt.Errorf("conectar ao banco de dados: %w", err)
	}
	defer func() {
		logger.Infof("fechando conexões com o banco de dados...")
		_ = db.Close()
	}()

	redisOptions, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("parse REDIS_URL: %w", err)
	}
	redisClient := redis.NewClient(redisOptions)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		_ = redisClient.Close()
		return fmt.Errorf("connect to Redis: %w", err)
	}
	defer redisClient.Close()
	workerCtx, stopWorkers := context.WithCancel(context.Background())
	defer stopWorkers()

	// 3. Configuração do Roteador e Servidor HTTP
	router := routes.SetupRouterWithContext(workerCtx, db, redisClient)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// 4. Inicia o servidor sem bloquear o processo principal
	serverErrors := make(chan error, 1)
	go func() {
		logger.Infof("servidor Gin iniciado em :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("iniciar servidor HTTP: %w", err)
		}
	}()

	// 5. Canal para escutar sinais de interrupção do SO (Docker stop, Ctrl+C)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Aguarda um erro de inicialização ou o sinal de desligamento
	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdown:
		logger.Infof("sinal de término recebido (%v), iniciando Graceful Shutdown...", sig)

		// Dá 10 segundos para requisições em andamento terminarem
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			// Força o encerramento se estourar o tempo limite
			_ = srv.Close()
			return fmt.Errorf("falha ao desligar servidor suavemente: %w", err)
		}
	}

	logger.Infof("servidor finalizado com sucesso")
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
