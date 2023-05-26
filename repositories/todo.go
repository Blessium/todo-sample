package repositories

import (
	old_error "errors"
	"github.com/blessium/todo-sample/database"
	"github.com/blessium/todo-sample/errors"
	"gorm.io/gorm"
	"time"
    "fmt"
)

var (
	TodoRepositoryPath = "repository/todo.go"
)

type ITodoRepository interface {
	AddTodo(todo Todo) (Todo, error)
	UpdateTodo(id uint, todo Todo) (Todo, error)
	UpdateTodos(todos []Todo) ([]Todo, error)
	GetTodo(id uint) (Todo, error)
	GetTodos() ([]Todo, error)
	DeleteTodo(id uint) error
	DeleteTodos() error
}

type Todo struct {
	ID           uint
	Title        string
	Notes        string
	CreationDate time.Time
	DueDate      time.Time
	Completed    bool
}

func (t Todo) mapToGorm() database.Todo {
	return database.Todo{
		ID:           t.ID,
		Title:        t.Title,
		Notes:        t.Notes,
		CreationDate: t.CreationDate,
		DueDate:      t.DueDate,
		Completed:    t.Completed,
	}
}

func mapFromGorm(t database.Todo) Todo {
	return Todo{
		ID:           t.ID,
		Title:        t.Title,
		Notes:        t.Notes,
		CreationDate: t.CreationDate,
		DueDate:      t.DueDate,
		Completed:    t.Completed,
	}
}

type GormTodoRepository struct {
	db *gorm.DB
}

func NewGormTodoRepository(db *gorm.DB) ITodoRepository {
	return GormTodoRepository{
		db: db,
	}
}

func newErrInternalDB(op string, details error) error {
    message := fmt.Sprintf("%s: internal database error", op)
    return errors.NewError(TodoRepositoryPath, errors.ErrInternal, message, details)
}

func (r GormTodoRepository) AddTodo(todo Todo) (Todo, error) {
	todo_db := todo.mapToGorm()

	result := r.db.Create(&todo_db)
	if result.Error != nil {
        return todo, newErrInternalDB("addTodo", result.Error)
	}

	return mapFromGorm(todo_db), nil
}

func (r GormTodoRepository) UpdateTodo(id uint, todo Todo) (Todo, error) {
	todo_db := todo.mapToGorm()

	temp, err := r.GetTodo(id)
    if err != nil {
		return Todo{}, err
	}

	todo_db.ID = id
    todo_db.CreationDate = temp.CreationDate

	result := r.db.Save(&todo_db)
	if result.Error != nil {
        return todo, newErrInternalDB("updateTodo", result.Error)
	}
    fmt.Print(todo_db)
	return mapFromGorm(todo_db), nil
}

func (r GormTodoRepository) UpdateTodos(todos []Todo) ([]Todo, error) {
	var todos_db []database.Todo
	for _, todo := range todos {
        temp, err := r.GetTodo(todo.ID)
        if err != nil {
            return nil, err
        }
        todo_to_add := todo.mapToGorm()
        todo_to_add.CreationDate = temp.CreationDate // ovveride creation_date
		todos_db = append(todos_db, todo_to_add)
	}

	result := r.db.Save(&todos_db)
	if result.Error != nil {
		return todos, newErrInternalDB("updateTodos", result.Error)
	}

	var todos_back []Todo
	for _, todo := range todos_db {
		todos_back = append(todos_back, mapFromGorm(todo))
        fmt.Print(todo)
	}

	return todos_back, nil
}

func (r GormTodoRepository) GetTodo(id uint) (Todo, error) {
	var todo_db database.Todo
	todo_db.ID = id

	result := r.db.Take(&todo_db)
	if result.Error != nil {
		if old_error.Is(result.Error, gorm.ErrRecordNotFound) {
			return Todo{}, errors.NewError(TodoRepositoryPath, errors.ErrNotExist, "todo doesn't exist", result.Error)
		}
		return Todo{}, newErrInternalDB("getTodo", result.Error)
	}

	return mapFromGorm(todo_db), nil
}

func (r GormTodoRepository) GetTodos() ([]Todo, error) {
	var todos_db []database.Todo

	result := r.db.Find(&todos_db)
	if result.Error != nil {
		return nil, newErrInternalDB("getTodos", result.Error)
	}

	var todos []Todo
	for _, todo_db := range todos_db {
		todos = append(todos, mapFromGorm(todo_db))
	}

	return todos, nil
}

func (r GormTodoRepository) DeleteTodo(id uint) error {
    _, err := r.GetTodo(id)
    if err != nil {
        return err
    }

	result := r.db.Delete(&database.Todo{ID: id})
	if result.Error != nil {
		return newErrInternalDB("deleteTodo", result.Error)
	}
	return nil
}

func (r GormTodoRepository) DeleteTodos() error {
	result := r.db.Where("1=1").Delete(&database.Todo{})
	if result.Error != nil {
		return newErrInternalDB("deleteTodos", result.Error)
	}
	return nil
}
