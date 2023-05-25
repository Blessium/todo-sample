package handlers

import (
	"github.com/blessium/todo-sample/services"
	"github.com/blessium/todo-sample/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
    "errors"
)

var (
    HandlerRequiredTitleError="todo schema: \"title\" field is required"
    HandlerFormatDueTimeError="todo schema: Title field is required"
)

type TodoHandler struct {
	todoService services.ITodoService
}

func NewTodoHandler(s services.ITodoService) TodoHandler {
	return TodoHandler{
		todoService: s,
	}
}

type TodoAddRequest struct {
	Title   string    `json:"title"`
	Notes   string    `json:"notes"`
	DueDate time.Time `json:"due_date"`
}

func (t TodoAddRequest) Validate() error {
    if t.Title == "" {
        return errors.New(HandlerRequiredTitleError)
    } 

    return nil
}

func (t TodoAddRequest) mapToService() services.Todo {
	return services.Todo{
		Title:   t.Title,
		Notes:   t.Notes,
		DueDate: t.DueDate,
	}
}

type TodoUpdateRequest struct {
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	DueDate   time.Time `json:"due_date"`
    Completed bool `json:"completed"`
}

func (t TodoUpdateRequest) mapToService() services.Todo {
	return services.Todo{
		Title:     t.Title,
		Notes:     t.Notes,
		DueDate:   t.DueDate,
		Completed: t.Completed,
	}
}

type TodoFullResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Notes        string    `json:"notes"`
	CreationDate time.Time `json:"creation_date"`
	DueDate      time.Time `json:"due_date"`
    Completed    bool `json:"completed"`
}

func FullResponseFromService(t services.Todo) TodoFullResponse {
	return TodoFullResponse{
		ID:           t.ID,
		Title:        t.Title,
		Notes:        t.Notes,
		CreationDate: t.CreationDate,
		DueDate:      t.DueDate,
		Completed:    t.Completed,
	}
}

func (t TodoHandler) AddTodo(c echo.Context) error {
	var r TodoAddRequest
	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

    if err := r.Validate(); err != nil {
        return c.String(http.StatusBadRequest, err.Error())
    }

	todo := r.mapToService()
	todo, err := t.todoService.AddTodo(todo)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, FullResponseFromService(todo))
}

func (t TodoHandler) DeleteTodo(c echo.Context) error {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := t.todoService.DeleteTodo(id); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (t TodoHandler) UpdateTodo(c echo.Context) error {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var r TodoUpdateRequest
	if err := c.Bind(&r); err != nil {
		return err
	}

	todo, err := t.todoService.UpdateTodo(id, r.mapToService())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, FullResponseFromService(todo))
}

func (t TodoHandler) GetTodo(c echo.Context) error {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	todo, err := t.todoService.GetTodo(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, FullResponseFromService(todo))
}

func (t TodoHandler) GetTodos(c echo.Context) error {
    todos, err := t.todoService.GetTodos()
    if err != nil {
        return c.String(http.StatusInternalServerError, err.Error())
    }
    var r_todos []TodoFullResponse
    for _, todo := range todos {
        r_todos = append(r_todos, FullResponseFromService(todo))
    }

    return c.JSON(http.StatusFound, r_todos)
}
