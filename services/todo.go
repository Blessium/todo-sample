package services

import (
	"errors"
	"github.com/blessium/todo-sample/repositories"
	"time"
)

var (
	InvalidTodoDateError = "todo: completition date is not valid"
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
		return errors.New(InvalidTodoDateError)
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

func (t TodoService) AddTodo(todo Todo) (Todo, error) {
	if err := todo.Validate(); err != nil {
		return todo, err
	}

	todo.CreationDate = time.Now()
	todo.Completed = false

	todo_repo := todo.mapToRepo()

    r_todo,  err := t.todoRepository.AddTodo(todo_repo); 
    if err != nil {
		return todo, err
	}

	return todoFromRepo(r_todo), nil
}
func (t TodoService) UpdateTodo(id uint, todo Todo) (Todo, error) {
	if err := todo.Validate(); err != nil {
		return todo, err
	}

	todo_repo := todo.mapToRepo()

	_, err := t.todoRepository.UpdateTodo(id, todo_repo)
    if err != nil {
		return todo, err
	}

	return todo, nil
}

func (t TodoService) UpdateTodos(todos []Todo) ([]Todo, error) {

	var todos_repo []repositories.Todo

	for _, todo := range todos {
		if err := todo.Validate(); err != nil {
			return nil, err
		} else {
			todos_repo = append(todos_repo, todo.mapToRepo())
		}
	}

	_, err := t.todoRepository.UpdateTodos(todos_repo)
    if err != nil {
		return nil, err
	}

	return todos, nil
}

func (t TodoService) GetTodo(id uint) (Todo, error) {
	todo_repo, err := t.todoRepository.GetTodo(id)
	if err != nil {
		return Todo{}, err
	}

	todo := todoFromRepo(todo_repo)

	return todo, nil
}

func (t TodoService) GetTodos() ([]Todo, error) {
	todos_repo, err := t.todoRepository.GetTodos()
	if err != nil {
		return nil, err
	}

	var todos []Todo

	for _, todo_repo := range todos_repo {
		todos = append(todos, todoFromRepo(todo_repo))
	}

	return todos, nil
}

func (t TodoService) DeleteTodo(id uint) error {
	if err := t.todoRepository.DeleteTodo(id); err != nil {
		return err
	}
	return nil
}

func (t TodoService) DeleteTodos() error {
	if err := t.todoRepository.DeleteTodos(); err != nil {
		return err
	}
	return nil
}
