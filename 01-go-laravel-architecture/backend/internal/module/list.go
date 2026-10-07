package module

import (
	"database/sql"

	"github.com/seu-usuario/taskflow-backend/internal/handler"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
	"github.com/seu-usuario/taskflow-backend/internal/service"
)

func BuildListModule(db *sql.DB) *handler.ListHandler {
	listRepository := repository.NewListRepository(db)
	listService := service.NewListService(listRepository)
	return handler.NewListHandler(listService)
}
