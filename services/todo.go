package services

import (
	"github.com/blessium/todo-sample/errors"
	"github.com/blessium/todo-sample/repositories"
	"time"
)

var (
	ServiceTodoPath = "service/todo.go"
)

type Todo struct {
	ID           uint
	Title        string
	Notes        string
	CreationDate time.Time
	DueDate      time.Time
	Completed    bool
}

func (t Todo) Validate() error {
	if t.DueDate.Unix() < time.Now().Unix() {
		return errors.NewError(ServiceTodoPath, errors.ErrValidation, "due_date is already expired", nil)
	}
	return nil
}

func (t Todo) mapToRepo() repositories.Todo {
	return repositories.Todo{
		ID:           t.ID,
		Title:        t.Title,
		Notes:        t.Notes,
		CreationDate: t.CreationDate,
		DueDate:      t.DueDate,
		Completed:    t.Completed,
	}
}

func todoFromRepo(t repositories.Todo) Todo {
	return Todo{
		ID:           t.ID,
		Title:        t.Title,
		Notes:        t.Notes,
		CreationDate: t.CreationDate,
		DueDate:      t.DueDate,
		Completed:    t.Completed,
	}
}

type TodoFilter struct {
	IsCompleted bool
}

type ITodoService interface {
	AddTodo(todo Todo) (Todo, error)
	UpdateTodo(id uint, todo Todo) (Todo, error)
	UpdateTodos(todos []Todo) ([]Todo, error)
	GetTodo(id uint) (Todo, error)
	GetTodos() ([]Todo, error)
	DeleteTodo(id uint) error
	DeleteTodos() error
}

type TodoService struct {
	todoRepository repositories.ITodoRepository
}

func NewTodoService(t repositories.ITodoRepository) ITodoService {
	return TodoService{
		todoRepository: t,
	}
}

func handleError(e error) error {
	internErr := errors.NewError(ServiceTodoPath, errors.ErrInternal, "internal server error", nil)
	// Non esporre i dettagli interni
	if errors.IsType(e, errors.ErrInternal) {
		return internErr
	} else if errors.IsType(e, errors.ErrValidation) || errors.IsType(e, errors.ErrNotExist) {
		return e
	}

	return nil
}

func (t TodoService) AddTodo(todo Todo) (Todo, error) {
	if err := todo.Validate(); err != nil {
		return todo, handleError(err)
	}

	todo.CreationDate = time.Now()
	todo.Completed = false

	todo_repo := todo.mapToRepo()

	r_todo, err := t.todoRepository.AddTodo(todo_repo)
	if err != nil {
		return todo, handleError(err)
	}

	return todoFromRepo(r_todo), nil
}
func (t TodoService) UpdateTodo(id uint, todo Todo) (Todo, error) {
	if err := todo.Validate(); err != nil {
		return todo, handleError(err)
	}

    if id == 0 {
		return Todo{}, errors.NewError(ServiceTodoPath, errors.ErrValidation, "id cannot be 0", nil)
    }

	todo_repo := todo.mapToRepo()

	res_todo, err := t.todoRepository.UpdateTodo(id, todo_repo)
	if err != nil {
		return todo, handleError(err)
	}

	return todoFromRepo(res_todo), nil
}

func (t TodoService) UpdateTodos(todos []Todo) ([]Todo, error) {

	var todos_repo []repositories.Todo

	for _, todo := range todos {
		if err := todo.Validate(); err != nil {
			return nil, handleError(err)
		} else {
			todos_repo = append(todos_repo, todo.mapToRepo())
		}
	}

	res_todos, err := t.todoRepository.UpdateTodos(todos_repo)
	if err != nil {
		return nil, handleError(err)
	}

    todos = todos[:0]
    for _, todo := range res_todos {
        todos = append(todos, todoFromRepo(todo))
    }

	return todos, nil
}

func (t TodoService) GetTodo(id uint) (Todo, error) {
	if id == 0 {
		return Todo{}, errors.NewError(ServiceTodoPath, errors.ErrValidation, "id cannot be 0", nil)
	}

	todo_repo, err := t.todoRepository.GetTodo(id)
	if err != nil {
		return Todo{}, handleError(err)
	}

	todo := todoFromRepo(todo_repo)

	return todo, nil
}

func (t TodoService) GetTodos() ([]Todo, error) {
	todos_repo, err := t.todoRepository.GetTodos()
	if err != nil {
		return nil, handleError(err)
	}

	var todos []Todo

	for _, todo_repo := range todos_repo {
		todos = append(todos, todoFromRepo(todo_repo))
	}

	return todos, nil
}

func (t TodoService) DeleteTodo(id uint) error {
	if id == 0 {
		return errors.NewError(ServiceTodoPath, errors.ErrValidation, "id cannot be 0", nil)
	}
	if err := t.todoRepository.DeleteTodo(id); err != nil {
		return handleError(err)
	}
	return nil
}

func (t TodoService) DeleteTodos() error {
	if err := t.todoRepository.DeleteTodos(); err != nil {
		return handleError(err)
	}
	return nil
}
