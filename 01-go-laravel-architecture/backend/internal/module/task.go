// Package module composes repository, service, and handler dependencies.
package module

import (
	"database/sql"

	"github.com/seu-usuario/taskflow-backend/internal/handler"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
	"github.com/seu-usuario/taskflow-backend/internal/service"
)

// BuildTaskModule constructs the task dependency graph from the shared database.
func BuildTaskModule(db *sql.DB) *handler.TaskHandler {
	taskRepository := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepository)
	return handler.NewTaskHandler(taskService)
}
