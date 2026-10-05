package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectMySQL lê as variáveis de ambiente e retorna um pool de conexões com o MySQL.
func ConnectMySQL() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// DSN (Data Source Name): formato padronizado do driver MySQL em Go
	// Ex: usuario:senha@tcp(db:3306)/taskflow_db?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão MySQL: %w", err)
	}

	// Configuração do Pool de Conexões (Equivalente às configurações de pool do config/database.php no Laravel)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Como o MySQL no Docker pode demorar alguns segundos para inicializar,
	// fazemos tentativas de Ping antes de falhar totalmente.
	var pingErr error
	for i := 1; i <= 10; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			log.Println("✅ Conexão com MySQL estabelecida com sucesso!")
			return db, nil
		}
		log.Printf("⏳ Aguardando MySQL iniciar... Tentativa %d/10\n", i)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("não foi possível conectar ao MySQL após 10 tentativas: %w", pingErr)
}
