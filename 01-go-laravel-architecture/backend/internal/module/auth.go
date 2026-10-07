package module

import (
	"database/sql"

	"github.com/seu-usuario/taskflow-backend/internal/handler"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
	"github.com/seu-usuario/taskflow-backend/internal/service"
)

func BuildAuthModule(db *sql.DB) *handler.AuthHandler {
	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	return handler.NewAuthHandler(authService)
}
