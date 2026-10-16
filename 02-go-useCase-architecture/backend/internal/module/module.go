package module

import (
	"database/sql"

	"github.com/seu-usuario/taskflow-backend/internal/handler"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
	authusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/auth"
	listusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/list"
	tagusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/tag"
	taskusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/task"
)

// AppModules exposes the handlers required by the HTTP router.
type AppModules struct {
	Auth  *handler.AuthHandler
	Lists *handler.ListHandler
	Tasks *handler.TaskHandler
	Tags  *handler.TagHandler
}

// Build creates the repository-to-handler dependency graph for the application.
func Build(db *sql.DB) AppModules {
	userRepository := repository.NewUserRepository(db)
	listRepository := repository.NewListRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	tagRepository := repository.NewTagRepository(db)

	register, login, getUser, updateProfile, updatePassword := authusecase.New(userRepository)
	createList := listusecase.NewCreateList(listRepository)
	getList := listusecase.NewGetListByID(listRepository)
	listLists := listusecase.NewListLists(listRepository)
	updateList := listusecase.NewUpdateList(listRepository)
	deleteList := listusecase.NewDeleteList(listRepository)
	createTask, updateTask, deleteTask, tasksByList, taskByID := taskusecase.New(taskRepository, listRepository)
	createTag, updateTag, deleteTag, tagsByList, tagByID := tagusecase.New(tagRepository, listRepository)

	return AppModules{
		Auth:  handler.NewAuthHandler(register, login, getUser, updateProfile, updatePassword),
		Lists: handler.NewListHandler(createList, getList, listLists, updateList, deleteList),
		Tasks: handler.NewTaskHandler(createTask, updateTask, deleteTask, tasksByList, taskByID),
		Tags:  handler.NewTagHandler(createTag, updateTag, deleteTag, tagsByList, tagByID),
	}
}
