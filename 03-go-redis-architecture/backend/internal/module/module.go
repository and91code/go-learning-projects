package module

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/seu-usuario/taskflow-backend/internal/handler"
	"github.com/seu-usuario/taskflow-backend/internal/pubsub"
	"github.com/seu-usuario/taskflow-backend/internal/queue"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
	authusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/auth"
	listusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/list"
	tagusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/tag"
	taskusecase "github.com/seu-usuario/taskflow-backend/internal/usecase/task"
)

// AppModules exposes the handlers required by the HTTP router.
type AppModules struct {
	Auth       *handler.AuthHandler
	Lists      *handler.ListHandler
	Tasks      *handler.TaskHandler
	Tags       *handler.TagHandler
	Queue      queue.Producer
	Dispatcher queue.Dispatcher
	Worker     interface{ Run(context.Context) }
}

// Build creates the repository-to-handler dependency graph for the application.
func Build(db *sql.DB) AppModules {
	return build(db, nil)
}
func BuildWithRedis(db *sql.DB, client *redis.Client) AppModules {
	return build(db, client)
}

func build(db *sql.DB, client *redis.Client) AppModules {
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

	modules := AppModules{
		Auth:  handler.NewAuthHandler(register, login, getUser, updateProfile, updatePassword),
		Lists: handler.NewListHandler(createList, getList, listLists, updateList, deleteList),
		Tasks: handler.NewTaskHandler(createTask, updateTask, deleteTask, tasksByList, taskByID),
		Tags:  handler.NewTagHandler(createTag, updateTag, deleteTag, tagsByList, tagByID),
	}
	if client != nil {
		modules.Queue = queue.NewProducer(client)
		modules.Lists.SetQueue(modules.Queue)
		modules.Tasks.SetQueue(modules.Queue)
		modules.Tags.SetQueue(modules.Queue)
		modules.Dispatcher = queue.DispatcherFuncs{
			CreateTask: createTask, UpdateTask: updateTask, DeleteTask: deleteTask,
			CreateList: createList, UpdateList: updateList, DeleteList: deleteList,
			CreateTag: createTag, UpdateTag: updateTag, DeleteTag: deleteTag,
		}
		modules.Worker = queue.NewWorker(client, modules.Dispatcher, pubsub.NewBroadcaster(client))
	}
	return modules
}
