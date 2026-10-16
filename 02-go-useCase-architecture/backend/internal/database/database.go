package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/seu-usuario/taskflow-backend/internal/logger"
)

const connectionTimeout = 5 * time.Second

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Connect(cfg Config) (*sql.DB, error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Name == "" {
		err := errors.New("DB_HOST, DB_PORT, DB_USER e DB_NAME devem estar configurados")
		logger.Errorf("configuração de conexão inválida: %v", err)
		return nil, err
	}

	driverConfig := mysql.Config{
		User:      cfg.User,
		Passwd:    cfg.Password,
		Net:       "tcp",
		Addr:      net.JoinHostPort(cfg.Host, cfg.Port),
		DBName:    cfg.Name,
		ParseTime: true,
	}

	db, err := sql.Open("mysql", driverConfig.FormatDSN())
	if err != nil {
		wrappedErr := fmt.Errorf("abrir conexão MySQL: %w", err)
		logger.Errorf("%v", wrappedErr)
		return nil, wrappedErr
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		wrappedErr := fmt.Errorf("verificar conexão MySQL: %w", err)
		logger.Errorf("%v", wrappedErr)
		return nil, wrappedErr
	}

	logger.Infof("conexão com MySQL estabelecida")
	return db, nil
}
