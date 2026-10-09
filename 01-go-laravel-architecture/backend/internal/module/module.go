package module

import (
	"database/sql"

	"github.com/seu-usuario/taskflow-backend/internal/handler"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
	"github.com/seu-usuario/taskflow-backend/internal/service"
)

// AppModules exposes the handlers required by the HTTP router.
type AppModules struct {
	Auth  *handler.AuthHandler
	Lists *handler.ListHandler
	Tasks *handler.TaskHandler
}

// Build creates the repository-to-handler dependency graph for the application.
func Build(db *sql.DB) AppModules {
	userRepository := repository.NewUserRepository(db)
	listRepository := repository.NewListRepository(db)
	taskRepository := repository.NewTaskRepository(db)

	authService := service.NewAuthService(userRepository)
	listService := service.NewListService(listRepository)
	taskService := service.NewTaskService(taskRepository, listRepository)

	return AppModules{
		Auth:  handler.NewAuthHandler(authService),
		Lists: handler.NewListHandler(listService),
		Tasks: handler.NewTaskHandler(taskService),
	}
}
