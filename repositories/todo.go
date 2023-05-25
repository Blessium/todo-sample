package repositories

import (
	"gorm.io/gorm"
    "github.com/blessium/todo-sample/database"
	"time"
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

func (r GormTodoRepository) AddTodo(todo Todo) (Todo, error) {
    todo_db := todo.mapToGorm()

    result := r.db.Create(&todo_db)
    if result.Error != nil {
        return todo, result.Error
    }

	return mapFromGorm(todo_db), nil
}

func (r GormTodoRepository) UpdateTodo(id uint, todo Todo) (Todo, error) {
    todo_db := todo.mapToGorm()

    result := r.db.Save(&todo_db)
    if result.Error != nil {
        return todo, result.Error
    }
	return mapFromGorm(todo_db), nil
}

func (r GormTodoRepository) UpdateTodos(todos []Todo) ([]Todo, error) {
    var todos_db []database.Todo
    for _, todo := range todos {
        todos_db = append(todos_db, todo.mapToGorm())
    }

    result := r.db.Save(&todos_db)
    if result.Error != nil {
        return todos, result.Error
    }

    var todos_back []Todo
    for _, todo := range todos_db {
        todos_back = append(todos_back, mapFromGorm(todo)) 
    }

	return todos_back, nil
}

func (r GormTodoRepository) GetTodo(id uint) (Todo, error) {
    var todo_db database.Todo
    todo_db.ID = id

    result := r.db.Take(&todo_db) 
    if result.Error != nil {
        return Todo{}, result.Error
    }

	return mapFromGorm(todo_db), nil
}

func (r GormTodoRepository) GetTodos() ([]Todo, error) {
    var todos_db []database.Todo

    result := r.db.Find(&todos_db)
    if result.Error != nil {
        return nil, result.Error
    }

    var todos []Todo
    for _, todo_db := range todos_db {
        todos = append(todos, mapFromGorm(todo_db))
    }


	return todos, nil
}

func (r GormTodoRepository) DeleteTodo(id uint) error {
    result := r.db.Delete(&database.Todo{ID : id}) 
    if result.Error != nil {
        return result.Error
    }
	return nil
}

func (r GormTodoRepository) DeleteTodos() error {
    result := r.db.Where("1=1").Delete(&database.Todo{})
    if result.Error != nil {
        return result.Error
    }
	return nil
}
